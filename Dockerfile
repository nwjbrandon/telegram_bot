FROM public.ecr.aws/lambda/provided:al2 as build

# Install build tools
RUN yum -y install git golang
RUN go env -w GOPROXY=direct

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download 

# Build
COPY main.go ./
RUN go build -o /main

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
ENTRYPOINT [ "/main" ] 