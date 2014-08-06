POSSIBLE += $(shell find db -type f | egrep -v '\.dat' | sed 's/$$/.dat/g')

all: ${POSSIBLE}

db/off/%.dat: db/off/%
	strfile -x $<

%.dat: %
	strfile $<

clean:
	rm -f db/*.dat db/off/*.dat
