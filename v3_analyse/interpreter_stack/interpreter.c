
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

unsigned char *memory;

int heapEnd = 0;

double current_start = 0;
int current_cc = 0;
struct analysis_info {
  double total_time;  /* sec diesmal */
  int64_t num_called;
};

struct analysis_info command_infos[31];

#ifdef DEBUG
    #define debug_print(fmt, args...) printf("%s:%04d: "fmt"\n", __FUNCTION__, __LINE__, ##args)
#else
    #define debug_print(fmt, args...)
#endif


#define COMMANDSIZE 17
/* Mit dispatch wird normal im code weitergegangen == cp um eins erhöht */
#define DISPATCH() pc = pc + COMMANDSIZE; /*printf("hab ccode: %d\n", code[pc]);*/ EXECUTE()
/* Mit execute wird dann wirklich das nächste kommando aufgerufen - wenn der pc manuell geändert wird, wird nur EXECUTE verwendet */
#define EXECUTE() ANA_WRITE(); debug_print("Now at position %d (%d) - command %d", (int)(pc / COMMANDSIZE), pc, code[pc]); ANA_START(); goto *command_impl[code[pc]]

#define ANA_START() clock_gettime(CLOCK_MONOTONIC, &spec); current_start = ((double)spec.tv_sec + 1.0e-9*spec.tv_nsec) * 1000; current_cc = code[pc];

#define ANA_WRITE() \
  if(current_cc != 0) { \
    clock_gettime(CLOCK_MONOTONIC, &spec); \
    command_infos[current_cc].total_time += (((double)spec.tv_sec + 1.0e-9*spec.tv_nsec) * 1000) - current_start; \
    command_infos[current_cc].num_called++; \
  } \

#define P1() ((int64_t)code[pc + 1 + 7] << 56 | (int64_t)code[pc + 1 + 6] << 48 | (int64_t)code[pc + 1 + 5] << 40 | (int64_t)code[pc + 1 + 4] << 32 | (int64_t)code[pc + 1 + 3] << 24 | (int64_t)code[pc + 1 + 2] << 16 | (int64_t)code[pc + 1 + 1] << 8 | (int64_t)code[pc + 1 + 0])
#define P2() ((int64_t)code[pc + 9 + 7] << 56 | (int64_t)code[pc + 9 + 6] << 48 | (int64_t)code[pc + 9 + 5] << 40 | (int64_t)code[pc + 9 + 4] << 32 | (int64_t)code[pc + 9 + 3] << 24 | (int64_t)code[pc + 9 + 2] << 16 | (int64_t)code[pc + 9 + 1] << 8 | (int64_t)code[pc + 9 + 0])
#define INT64_AT(buffer) ((int64_t)buffer[7] << 56 | (int64_t)buffer[6] << 48 | (int64_t)buffer[5] << 40 | (int64_t)buffer[4] << 32 | (int64_t)buffer[3] << 24 | (int64_t)buffer[2] << 16 | (int64_t)buffer[1] << 8 | (int64_t)buffer[0])
#define INT_AT(buffer) ((int)(INT64_AT(buffer)))

#define RETURN(retval) free(stack); return retval;

int interpret(unsigned char *code) {
  setbuf(stdout, NULL);
  int64_t *stack;
  /* Stackpointer zeigt immer auf die Speicherzelle in die als nächstes geschrieben werden kann */
  /* Push ist also schreiben und dann erhöhen */
  /* Pop ist vermindern und dann lesen */
  int stackpointer = 0;
  int framepointer = 0;

  int i = 0;
  int anzparam;

  for(int i=0;i<31;i++) {
    command_infos[i].total_time = 0;
    command_infos[i].num_called = 0;
  }

  struct timespec spec;


  unsigned char tempbyte;
  /* Hier stehen die command-implementierungen */
  static void* command_impl[] = {
    &&op_invalid, // 0
    &&op_push,
    &&op_store,
    &&op_retrieve,
    &&op_echos,
    &&op_progentry,
    &&op_progexit,
    NULL,
    NULL,
    &&op_add,
    &&op_mult, // 10
    &&op_subt,
    &&op_div,
    &&op_echoi,
    &&op_concat,
    &&op_gt,
    &&op_lt,
    &&op_neq,
    &&op_eq,
    &&op_jfalse,
    &&op_j, // 20
    &&op_alloc,
    &&op_storei,
    &&op_retrievei,
    &&op_len,
    &&op_disc,
    &&op_call,
    &&op_funcentry,
    &&op_funcexit,
    NULL,
    &&op_return // 30
  };

  stack = (int64_t*)malloc(5000 * sizeof(int64_t));
  if(stack == NULL) {
    printf("alloc stack failed");
    return 1;
  }

  int pc = 0 - COMMANDSIZE;

  DISPATCH();

  op_invalid:
    debug_print("Invalid command 0");
    RETURN(1);
  op_progentry:
    stackpointer = framepointer + P1();
    DISPATCH();
  op_push:
    stack[stackpointer] = P1();
    stackpointer++;
    DISPATCH();
  op_store:
    stackpointer--;
    stack[framepointer+(int)P1()] = stack[stackpointer];
    DISPATCH();
  op_retrieve:
    stack[stackpointer] = stack[framepointer + (int)P1()];
    stackpointer++;
    DISPATCH();
  op_echos:
    stackpointer--;

    //printf("echos of string size: %ld\n", INT64_AT((memory + stack[stackpointer])));
    // das byte eins nach dem string wird ausgelagert und dort hin der null delimiter geschrieben
    tempbyte = memory[stack[stackpointer] + 8 + /* length */INT_AT((memory + stack[stackpointer]))];

    // NULL byte
    memory[stack[stackpointer] + 8 + /* length */INT_AT((memory + stack[stackpointer]))] = 0;

    fputs(memory + stack[stackpointer] + 8, stdout);

    // tempbyte zurückschreiben
    memory[stack[stackpointer] + 8 + /* length */INT_AT((memory + stack[stackpointer]))] = tempbyte;
    DISPATCH();
  op_progexit:
    RETURN(0);
  op_add:
    stack[stackpointer-2] = stack[stackpointer-2] + stack[stackpointer-1];

    /* der stack ist jetzt eins weniger */
    stackpointer--;
    DISPATCH();
  op_mult:
    stack[stackpointer-2] = stack[stackpointer-2] * stack[stackpointer-1];

    /* der stack ist jetzt eins weniger */
    stackpointer--;
    DISPATCH();
  op_subt:
    stack[stackpointer-2] = stack[stackpointer-2] - stack[stackpointer-1];

    /* der stack ist jetzt eins weniger */
    stackpointer--;
    DISPATCH();
  op_div:
    stack[stackpointer-2] = stack[stackpointer-2] / stack[stackpointer-1];

    /* der stack ist jetzt eins weniger */
    stackpointer--;
    DISPATCH();
  op_echoi:
    stackpointer--;
    fprintf(stdout, "%ld", stack[stackpointer]);
    DISPATCH();
  op_concat:
    RETURN(1);
  op_gt:
    stack[stackpointer-2] = (int64_t)(stack[stackpointer-2] > stack[stackpointer-1]);
    debug_print("gt: %d", (int)stack[stackpointer-2]);
    stackpointer -= 1;
    DISPATCH();
  op_lt:
    /* 2 werte am stack */
    /* kleiner vergleich - ergebnis am stack*/
    stack[stackpointer-2] = (int64_t)(stack[stackpointer-2] < stack[stackpointer-1]);
    debug_print("lt: %d", (int)stack[stackpointer-2]);
    stackpointer -= 1;
    DISPATCH();
  op_neq:
    stack[stackpointer-2] = (int64_t)(stack[stackpointer-2] != stack[stackpointer-1]);
    debug_print("neq: %d", (int)stack[stackpointer-2]);
    stackpointer -= 1;
    DISPATCH();
  op_eq:
    stack[stackpointer-2] = (int64_t)(stack[stackpointer-2] == stack[stackpointer-1]);
    debug_print("eq: %d", (int)stack[stackpointer-2]);
    stackpointer -= 1;
    DISPATCH();
  op_jfalse:
    /* Jump wenn wert am stackpointer == 0 */
    stackpointer--;
    if(stack[stackpointer] == 0) {
      pc = pc + (COMMANDSIZE * ((int)P1()));
      EXECUTE();
    } else {
      DISPATCH();
    }
  op_j:
    /* relative jump + P1() */
    pc = pc + (COMMANDSIZE * ((int)P1()));
    EXECUTE();
  op_alloc: RETURN(1);
  op_storei: RETURN(1);
  op_retrievei: RETURN(1);
  op_len: RETURN(1);
  op_disc:
    stackpointer--;
    DISPATCH();
  op_call:
    anzparam = P2();
    /* Aufrufparameter verschieben */
    for (i = 0; i < anzparam; i++) {
      stackpointer--;
      stack[stackpointer+3] = stack[stackpointer];
    }
    // steht jetzt auf retv
    stack[stackpointer] = 0; /*standard return value = 0 */
		stackpointer++;
		// SP steht jetzt auf zelle für ofp
		stack[stackpointer] = framepointer;
		stackpointer++;
		// SP steht jetzt auf zelle für oriip
		stack[stackpointer] = pc + COMMANDSIZE; /* es wird +1 gerechnet, damit er beim return nur übernommen werden muss */
		stackpointer++;
    debug_print("calling function (retv: %d, ofp: %d, oip: %d)", (int)stack[stackpointer-3], (int)stack[stackpointer-2], (int)(stack[stackpointer-1]/17));

    stackpointer += anzparam;
    pc = pc + (COMMANDSIZE * P1());
    EXECUTE();
  op_funcentry:
    /* setzt nur noch den framepointer richtig */
    stackpointer = stackpointer + P2();
    framepointer = stackpointer - P1();

    debug_print("new fp: %d", framepointer);
    debug_print("new sp: %d", stackpointer);

    DISPATCH();
  op_funcexit:
    debug_print("fp before return: %d", framepointer);
    debug_print("Returning (retv: %d, ofp: %d, oip: %d)", (int)stack[framepointer-3], (int)stack[framepointer-2], (int)(stack[framepointer-1]/17));
    stackpointer = framepointer - 2; /* -2 damit er ofp und orip überspringt */
  	/* steht jetzt richtig (nach retval) */
  	framepointer = stack[stackpointer]; /* wert nach retval = framePointer */


    pc = stack[stackpointer+1];
    debug_print("After return fp: %d, sp: %d", framepointer, stackpointer);
    EXECUTE();
  op_return:
		/* return value steht 3 vor framepointer */
		stackpointer--;
		stack[framepointer-3] = stack[stackpointer];
    goto op_funcexit; // Weitermachen mit funcexit
}


int main(int argc, char* argv[])
{
  memory = malloc(3000);

  unsigned char *buffer;
  FILE *f = fopen(argv[1], "rb");
  if ( f != NULL )
  {
    fseek(f, 0L, SEEK_END);
    long s = ftell(f);  // ENDE
    rewind(f); // wieder anfangen
    buffer = malloc(s);
    if ( buffer != NULL )
    {
      int readsize = fread(buffer, s, 1, f);
      // we can now close the file
      fclose(f);
      f = NULL;

      debug_print("read program with size: %d\n", (int)(readsize*s));

      int64_t litsize = INT64_AT(buffer);

      debug_print("literal size: %d\n", (int)litsize);

      memcpy(memory, buffer + 8, litsize);

      heapEnd = litsize;

      unsigned char *buffer_commands = buffer + 8 + litsize;

      // Habe jetzt buffer
      // Literals wegspeichern

      interpret(buffer_commands);

      printf("\n");
      for(int i=0;i<31;i++) {
        if(command_infos[i].num_called > 0)
          printf("%d;%lf;%lf;%ld\n", i, command_infos[i].total_time, command_infos[i].total_time/(command_infos[i].num_called), command_infos[i].num_called);
      }

      free(buffer);
    }
    if (f != NULL) fclose(f);
  }
  return EXIT_SUCCESS;
}
