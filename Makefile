test: jackslinks.test

jackslinks.test: *.go
	go test -x -v -cpuprofile=profile.out

profile.out: jackslinks.test
	./jackslinks.test -test.cpuprofile=profile.out -test.bench=.

profile: profile.out
	go tool pprof -web jackslinks.test profile.out

clean:
	rm -rf profile.out jackslinks.test
