
## Run Step
1.  cd $PROJ_ROOT
2.  go mod vendor
3.  go build -mod=vendor -o ./main ./src/apps/linting_app/linting_app.go
4.  chmod +x ./main
5.  ./main (be aware of user should have permission of folder)
6.  chmod +x ./test-case.sh && sh ./test-case.sh (use the bash you have)