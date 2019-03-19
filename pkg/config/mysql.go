package config

import (
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type (
	MySqlConfig struct {
		file    string
		Current string             `yaml:"current"`
		Configs []*mySqlConfigItem `json:"configs"`
	}

	mySqlConfigItem struct {
		Name string `yaml:"name"`
		Addr string `yaml:"addr"`
	}
)

func GetMysqlConfig(file string) (*MySqlConfig, error) {
	var c MySqlConfig
	c.file = file
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return &c, err
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return &c, err
	}

	return &c, nil
}

func (t *MySqlConfig) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "addr"})
	for _, c := range t.Configs {
		name := c.Name
		if c.Name == t.Current {
			name = "*" + c.Name
		}
		table.Append([]string{name, c.Addr})
	}

	table.Render()
}

func (t *MySqlConfig) Set(name, addr string) {
	var item *mySqlConfigItem
	for _, c := range t.Configs {
		if c.Name == name {
			item = c
			c.Addr = addr
			break
		}
	}

	if item == nil {
		item = &mySqlConfigItem{
			Name: name,
			Addr: addr,
		}
		t.Configs = append(t.Configs, item)
		if len(t.Configs) == 1 {
			t.Current = item.Name
		}
	}

	t.save()
}

func (t *MySqlConfig) Remove(name string) {
	pos := -1
	for i, c := range t.Configs {
		if c.Name == name {
			pos = i
			break
		}
	}

	if pos != -1 {
		t.Configs = append(t.Configs[:pos], t.Configs[pos+1:]...)
	}

	if len(t.Configs) == 1 {
		t.Current = t.Configs[0].Name
	}

	t.save()
}

func (t *MySqlConfig) Use(name string) {
	for _, c := range t.Configs {
		if c.Name == name {
			t.Current = name
			break
		}
	}
	t.save()
}

func (t *MySqlConfig) CurrentAddr() string {
	for _, c := range t.Configs {
		if c.Name == t.Current {
			return c.Addr
		}
	}
	return ""
}

func (t *MySqlConfig) save() {
	b, err := yaml.Marshal(t)
	if err != nil {
		panic(err)
	}

	_ = ioutil.WriteFile(t.file, b, 0600)
}
