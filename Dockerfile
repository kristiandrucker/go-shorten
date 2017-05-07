FROM scratch
COPY go-shorten templates .
COPY templates templates
ENTRYPOINT ["./go-shorten"]