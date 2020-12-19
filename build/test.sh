if [[ -z $CI_PROJECT_DIR ]]
then
    echo "Missing CI_PROJECT_DIR"
    exit 1
fi
cd $CI_PROJECT_DIR/src
go test -v ./internal/...