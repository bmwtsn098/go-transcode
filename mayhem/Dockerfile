FROM golang:1.18 as builder

COPY . /go-transcode
WORKDIR /go-transcode/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o fuzzHlsproxy.a && \
    clang -fsanitize=fuzzer fuzzHlsproxy.a -o fuzzHlsproxy.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /go-transcode/mayhem/fuzzHlsproxy.libfuzzer /