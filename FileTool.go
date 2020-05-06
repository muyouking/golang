package Tools

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)


//
func Getcwd() string {
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir,err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	return  dir
}
//直接获取指定文件夹的某种类型的文件

func H_GetFileList(pathdir,ext string) []string{
	if !IsDir(pathdir){
		return nil
	}
	//20200506扩展,如果exe为""空，则表示返回这个文件夹中的所有文件
	if ext == ""{
		ext = "*"
	}
	extfilelist,err := filepath.Glob(path.Join(pathdir,"*." + ext))
	if err != nil{
		fmt.Println(err)
	}
	return extfilelist

}
func H_GetFileContents(path string) (string,error)  {
	contentByte ,err:= ioutil.ReadFile(path)
	content := string(contentByte)
	//fmt.Println(content)
	return content,err
}
func LoadINI(name string)(string,[]string){
	cfg,err := ini.Load(name)
	if err != nil{
		log.Fatal("Failed to load:",err)
	}
	PA := cfg.Section("Section").Key("PATHA").String()
	PB := cfg.Section("Section").Key("PATHB").String()
	//fmt.Println("PathA:",PA)
	//fmt.Println("PathB:",PB)
	var pathB =make([]string,4)
	if strings.Contains(PB,","){
		pathB = strings.Split(PB,",")
		//fmt.Println(pathB)
	}
	return  PA,pathB

}
func Cmd (s string){
	cmd := exec.Command("cmd.exe","/c", " " + s)
	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}
}
func Cmd2 (s []string){
	cmd := exec.Command("cmd.exe",s...)
	err := cmd.Run()
	if err != nil{
		log.Fatal(err)
	}
}
//创建channel
func source(files []string) <-chan string {
	out := make(chan string, len(files))
	go func() {
		for _,file := range files{
			out <- file
		}
		close(out)
	}()
	return out
}
func WriteFiles(name string,s []string) error  {
	f,err := os.OpenFile(name,os.O_CREATE|os.O_TRUNC,0777)
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()


	w := bufio.NewWriter(f)

	for _,s1 := range s{
		_,err := w.WriteString(s1)
		if err != nil{
			log.Fatal(err)
		}
	}
	w.Flush()
	return  nil
}
//floderA: 原文件夹
//floderB：目标文件夹
//检查两个文件夹的文件是否一致
func CheckABfile(floderA,floderB string)([]string,error)   {
	floderA_fileInfo := GetFileNameAndModificationDate(floderA)
	floderB_fileInfo := GetFileNameAndModificationDate(floderB)
	recordfiles := make([]string,64)
		for filename,_ := range floderA_fileInfo{
			_,OK := floderB_fileInfo[filename]
			if OK{
				//如果大于10则记录下来，这是一个需要重新复制的文件
				if floderA_fileInfo[filename] - floderB_fileInfo[filename] < 10{
					continue
				}
				//如果小于10，则代表这个是不需要处理的文件
				s ,_ := recordfile(filepath.Join(floderA,filename),floderB)
				recordfiles = append(recordfiles,s)
			} else {
				s ,_ := recordfile(filepath.Join(floderA,filename),floderB)
				recordfiles = append(recordfiles,s)
			}
			//如果OK不为真，则代表这个文件在B文件夹里不存在，是新文件，需要复制过去



		}
		return recordfiles, nil
}
//一对多
func CheckABCFiles(f string,fs ...string)([]string)  {
	recordfiles := make([]string,64)
	for _,fi := range fs{
		NeedFile,err := CheckABfile(f,fi)
		if err != nil{
			log.Fatal(err)
		}
		recordfiles=append(recordfiles, NeedFile...)
	}
	return  recordfiles


}
//提取文件的复制命令
func recordfile(s string, s2 string) (string,error) {
	if !IsFile(s) {
		return "",errors.New(s + "不是文件")
	}
	if !IsDir(s2){
		return "",errors.New(s2 + "不是文件夹")
	}
	return ("copy /y " + s + " " + s2 +"\n"),nil
}
//获取文件的名字与修改日期
func GetFileNameAndModificationDate(path string) map[string] int64 {
	filesA,err := WalkDir(path,"")
	if err != nil{
		log.Fatal(err)
	}
	filemap := make(map[string] int64)
	for _,file := range filesA{
		//fmt.Println(file)
		fileInfomation,err := GetFileInfo(file,false)
		if err != nil{
			log.Fatal(err)
		}
		//fmt.Println(fileInfomation)
		modtime := fileInfomation.filemodtime

		filename := fileInfomation.filename

		filemap[filename]=modtime
	}
	return filemap
}
/* 获取指定路径下以及所有子目录下的所有文件，可匹配后缀过滤（suffix为空则不过滤）*/
func WalkDir(dir, suffix string) (files []string, err error) {
	files = []string{}
	err = filepath.Walk(dir,func (fname string, fi os.FileInfo, err error) error {
		
		if fi.IsDir() {
			//忽略目录
			return nil
		}

		if len(suffix) == 0 || strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
			//文件后缀匹配
			files = append(files, fname)
		}

		return nil
	})	//filepath.Walk，第二个参数结束（函数作为参数）

	return files, err
}

//获取文件修改时间 返回unix时间戳
func GetModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

//判断文件夹是否存在
func IsDir(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return info.IsDir()
	}
	return false
}
func IsFile(name string) bool  {
	return !IsDir(name)
}
//检查文件是否存在
func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

//创建文件夹（如果文件夹不存在则创建）
func MakeDir(dir string) error {
	if !FileIsExisted(dir) {
		if err := os.MkdirAll(dir, 0777); err != nil { //os.ModePerm
			fmt.Println("MakeDir failed:", err)
			return err
		}
	}
	return nil
}

//复制文件
func Copyfile(orgfile, dstfile string) (int, error) {
	fn1, err := os.Open(orgfile)
	if err != nil {
		return 0, err
	}
	fn2, err := os.OpenFile(dstfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer fn1.Close()
	defer fn2.Close()
	//读写
	bs := make([]byte, 1024, 1024)
	n := -1 //读取的数据量
	total := 0
	for {
		n, err = fn1.Read(bs)
		if err == io.EOF || n == 0 {
			fmt.Println("拷贝完毕...")
			break
		} else if err != nil {
			fmt.Println("报错了!!!")
			return total, err
		}
		total += n
		fn2.Write(bs[:n])
	}
	return total, nil

}

func CopyFile2(src, des string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	//获取源文件的权限
	fi, _ := srcFile.Stat()
	perm := fi.Mode()

	//desFile, err := os.Create(des)  //无法复制源文件的所有权限
	desFile, err := os.OpenFile(des, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm) //复制源文件的所有权限
	if err != nil {
		return 0, err
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}
//判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
//获取文件相关信息
//filename:文件路径
//fileinfo:文件信息(路径名，文件名，md5,修改时间)
type FileInfo struct {
	fileDir string
	filename string
	filemd5 string
	filemodtime int64
}
func GetFileInfo(filename string,haveMd5 bool) (fileinfo *FileInfo, err  error){
	fileinfo = &FileInfo{} //路径名，文件名，MD5 ,修改时间
	//判断文件是否存在
	_ ,err = os.Stat(filename)
	if err != nil{
		fmt.Println("os.Stat err = ",err)
		return fileinfo,err
	}
	//提取信息
	fileinfo.fileDir,fileinfo.filename = filepath.Split(filename)	//文件路径与文件名
	fileinfo.filemodtime = GetModTime(filename)	//文件的修改时间
	//是否需要Md5
	HaveMd5:=haveMd5
	if HaveMd5==true{
		fileinfo.filemd5 = Md5sum3(filename)
		}

	
	return fileinfo,nil
}
//获取目标路径里的jpg,png,jpeg的图片序列
func H_GetPiclist(path string) []string {
	jpglist := H_GetFileList(path,"jpg")
	pnglist := H_GetFileList(path,"png")
	jpeglist := H_GetFileList(path,"jpeg")


	piclist := SliceJoin(jpglist,pnglist,jpeglist)
	fmt.Println("jpg的数量为:",len(jpglist))
	fmt.Println("png:",len(pnglist))
	fmt.Println("jpeg的数量为:",len(jpeglist))
	fmt.Println("总数为:",len(piclist))
	return piclist
}
///GO 并发
///filelist 待处理的切片
///f 具体执行的函数 --例如，遍历filelist,拿到每一个元素，作为参加数传给f函数进行处理

func H_Goroutine(filelist []string, f func(s string, c chan string)) {
	c := make(chan string, len(filelist))
	for _, v := range filelist {
		go f(v, c) //do something

	}
	for i := 0; i < len(filelist); i++ {
		<-c
	}
	close(c)
}

//具体的执行方法
func doSomething(s string, c chan string) {

	fmt.Println(s)
	c <- s
}