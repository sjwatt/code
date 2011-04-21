include $(GOROOT)/src/Make.inc

.SUFFIXES : .$O .go

.go.$O:
	$(GC) -o $@ $<

TARG = DHClient

OBJS = $(TARG).$O 

$(TARG): $(OBJS)
	$(LD) -o $@ $@.$O
	@echo "Executable is: $@"

clean:
	rm -rf *.[$(OS)o]

nuke:
	rm -rf *.[$(OS)o] *~ $(TARG)

test: nuke $(TARG)
	./$(TARG)

all:
