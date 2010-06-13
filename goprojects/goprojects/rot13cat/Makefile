include $(GOROOT)/src/Make.$(GOARCH)

.SUFFIXES : .$O .go

.go.$O:
	$(GC) -o $@ $<

TARG = r13cat

OBJS = file.$O $(TARG).$O 

$(TARG): $(OBJS)
	$(LD) -o $@ $@.$O
	@echo "Executable is: $@"

clean:
	rm -rf *.[$(OS)o]

nuke:
	rm -rf *.[$(OS)o] *~ $(TARG)

test: nuke $(TARG)
	./$(TARG) -rot13

all:
