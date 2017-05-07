FROM scratch
COPY go-shorten .
COPY templates templates
CMD ["./go-shorten"]
