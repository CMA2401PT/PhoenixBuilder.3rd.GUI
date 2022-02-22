#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>
#include <string.h>

// Decided to use --wrap described by GNU's ld(1) at first, but it seems that darwin's 
// ld didn't implement it, so uses objcopy in Makefile instead.

char args_isDebugMode=0;
char args_disableHashCheck=0;
char replaced_auth_server=0;
char *newAuthServer;
char args_muteWorldChat=0;
char args_noPyRpc=0;
char args_noNBT=0;

int main(int argc, char **argv, const char **envp);

void print_help(const char *self_name) {
	printf("%s [options]\n",self_name);
	printf("\t--debug: Run in debug mode.\n");
	printf("\t-A <url>, --auth-server=<url>: Use the specified authentication server, instead of the default one.\n");
	printf("\t--no-hash-check: Disable the hash check.\n");
	printf("\t-M, --no-world-chat: Ignore world chat on client side.\n");
	printf("\t--no-pyrpc: Disable the PyRpcPacket interaction, the client's commands will be prevented from execution by netease's rental server.\n");
	printf("\t--no-nbt: Disable NBT Construction feature.\n");
	printf("\n");
	printf("\t-h, --help: Show this help context.\n");
}



int __wrap_main(int argc, char **argv, const char **envp) {
	while(1) {
		static struct option opts[]={
			{"debug", no_argument, 0, 0}, // 0
			{"help", no_argument, 0, 'h'}, // 1
			{"auth-server", required_argument, 0, 'A'}, //2
			{"no-hash-check", no_argument, 0, 0}, //3
			{"no-world-chat", no_argument, 0, 'M'}, //4
			{"no-pyrpc", no_argument, 0, 0}, //5
			{"no-nbt", no_argument, 0, 0}, //6
			{0, 0, 0, 0}
		};
		int option_index;
		int c=getopt_long(argc,argv,"hA:", opts, &option_index);
		if(c==-1)
			break;
		switch(c) {
		case 0:
			switch(option_index) {
			case 0:
				args_isDebugMode=1;
				break;
			case 3:
				args_disableHashCheck=1;
				break;
			case 5:
				args_noPyRpc=1;
				break;
			case 6:
				args_noNBT=1;
				break;
			};
			break;
		case 'h':
			print_help(argv[0]);
			return 0;
		case 'A':
			replaced_auth_server=1;
			size_t loo=strlen(optarg);
			newAuthServer=malloc(loo+1);
			memcpy(newAuthServer,optarg,loo+1);
			break;
		case 'M':
			args_muteWorldChat=1;
			break;
		default:
			print_help(argv[0]);
			return 1;
		};
	};
	return main(argc,argv,envp);
}