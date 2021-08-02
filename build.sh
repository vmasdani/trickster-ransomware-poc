rm -rf toencrypt/* &&\
rm -rf dist &&\
# xgo --targets=linux/amd64 . &&\
go build &&\
zip -r release.zip . &&\
mkdir dist &&\
mv release.zip dist
