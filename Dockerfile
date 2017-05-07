FROM scratch
COPY go-shorten .
CMD ["./go-shorten"]