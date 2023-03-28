CC := gcc.exe

ifeq ($(OS), Windows_NT)
	TARGET := delegate.lib
	SOURCES := $(wildcard win32*.c)
	RM := del
	RMARGS := /Q
else
	TARGET := delegate.a
	SOURCES := $(cocoa*.m)
	RM := rm
	RMARGS := -rf
endif

OBJS := $(SOURCES:%.c=%.o)

LDFLAGS := -no-pie

all: $(TARGET)
	
clean:
# @echo $(OS)
	$(RM) $(RMARGS) *.o
	$(RM) *.a
	$(RM) *.lib

$(TARGET): $(OBJS)
	ar rcs $@ $<

.o: .c
	$(CC) -o $@ -c $<
	