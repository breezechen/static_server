#Makefile
PROGS=server

all:$(PROGS)

CC=g++

OBJS= connection.o \
	custom_handler.o \
	html.o \
	io_service_pool.o \
	main.o \
	mime_types.o \
	reply.o \
	request_handler.o \
	request_parser.o \
	server.o

CFLAGS=-pthread -O2 -DNDEBUG

LIBS=-lboost_system -lboost_filesystem -lboost_thread -lrt

.cpp.o:
	$(CC) $(CFLAGS) -c $<

server:$(OBJS)
	$(CC) -o ../_builds/$@ $(OBJS) $(LIBS)


