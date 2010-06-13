include $(GOROOT)/src/Make.$(GOARCH)

.SUFFIXES : .$O .go

.go.$O:
	$(GC) -o $@ $<

TARG = puzzle

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
