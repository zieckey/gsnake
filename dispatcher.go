package gsnake

import (
	"github.com/golang/glog"
	"github.com/howeyc/fsnotify"
	"path/filepath"
	"sync"
	"strings"
)

type Dispatcher struct {
	dir        string
	watcher    *fsnotify.Watcher
	status     *ProcessStatus
	h          *FilesHandler
	textModule TextModule
	pcapModule PcapModule
}

var dispatcher *Dispatcher

func NewDispatcher() (*Dispatcher, error) {
	var err error
	dispatcher, err = newDispatcher(*dir)
	return dispatcher, err
}

func newDispatcher(dir string) (d *Dispatcher, err error) {
	glog.Infof("NewDispatcher")
	d = &Dispatcher{}
	d.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		glog.Fatal(err)
	}

	d.dir = dir
	d.status, err = NewProcessStatus(*statusFile)
	if err != nil || d.status == nil {
		glog.Fatal(err)
	}

	d.h, err = NewFilesHandler(dir)
	if err != nil {
		glog.Fatal(err)
	}

	d.textModule = &DefaultTextModule{}
	d.pcapModule = &DefaultPcapModule{}
	return d, err
}

func (d *Dispatcher) onCreate(ev *fsnotify.FileEvent) {
	if IsDir(ev.Name) {
		d.watcher.Watch(ev.Name)
		//Ignore this : FIXME if we renamed ev.Name latterly, we should add the new name to the watching list.
	} else {
		if ok, _ := filepath.Match(*filePattern, filepath.Base(ev.Name)); ok {
			d.h.OnFileCreated(ev.Name)
		} else {
			glog.Infof("Create a file <%v> but does not match the file pattern <%v>", ev.Name, *filePattern)
		}
	}
}

func (d *Dispatcher) onDelete(ev *fsnotify.FileEvent) {
	d.status.OnFileDeleted(ev.Name)
}

func (d *Dispatcher) onModify(ev *fsnotify.FileEvent) {
	d.h.OnFileModified(ev.Name)
}

func (d *Dispatcher) watchEvent(wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case ev := <-d.watcher.Event:
			if strings.ToLower(ev.Name) != strings.ToLower(*statusFile) {
				glog.Info("event:", ev, " name=", ev.Name)
				if ev.IsCreate() {
					d.onCreate(ev)
				} else if ev.IsDelete() {
					d.onDelete(ev)
				} else if ev.IsModify() {
					d.onModify(ev)
				} else {
					glog.Info("don't care this event:", ev)
				}
			}
		case err := <-d.watcher.Error:
			glog.Info("error:", err)
		}
	}
}

func (d *Dispatcher) Run() {
	glog.Infof("Watching <%v>", d.dir)
	err := d.watcher.Watch(d.dir)
	if err != nil {
		glog.Fatal("Watch event of " + d.dir + " FAILED: " + err.Error())
	}

	//start to watch the file event and wait the goroutine started
	var wg sync.WaitGroup
	wg.Add(1)
	go d.watchEvent(&wg)
	wg.Wait()

	//start file handler to run
	d.h.Run()

	d.Close()
}

func (d *Dispatcher) Close() {
	d.watcher.Close()
}

func (d *Dispatcher) RegisterTextModule(m TextModule) {
	d.textModule = m
}

func (d *Dispatcher) RegisterPcapModule(m PcapModule) {
	d.pcapModule = m
}
