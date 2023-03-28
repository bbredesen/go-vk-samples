ifeq ($(OS), Windows_NT)
	CC := gcc.exe
	TARGET := delegate.lib
	CSOURCES := $(wildcard win32*.c)
	RM := del
	RMARGS := /Q
else
	CC := gcc
	TARGET := libdelegate.a
	MSOURCES := $(wildcard cocoa*.m)
	RM := rm
	RMARGS := -rf
endif


COBJS := $(CSOURCES:%.c=%.o)
MOBJS := $(MSOURCES:%.m=%.o)
OBJS := $(COBJS) $(MOBJS)

LDFLAGS := -no-pie

all: $(TARGET)
	
clean:
# @echo $(OS)
	$(RM) $(RMARGS) *.o
	$(RM) $(RMARGS) $(TARGET)

$(TARGET): $(OBJS)
	ar rcs $@ $(OBJS)

.o: .c .m
	$(CC) -o $@ -c $<
	