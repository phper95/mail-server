package cmd

import (
	"fmt"
	"github.com/phper95/mail-server/internal/app/router"
	"github.com/phper95/mail-server/internal/conf"
	"github.com/phper95/mail-server/internal/tools/dkim"
	"github.com/urfave/cli"
)

var Dkim = cli.Command{
	Name:        "dkim",
	Usage:       "This command make dkim config file",
	Description: `Configure domain name settings`,
	Action:      makeDkim,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "Custom configuration file path"),
	},
}

func makeDkim(c *cli.Context) error {

	err := router.GlobalInit(c.String("config"))
	if err != nil {
		return err
	}

	domain := conf.Web.Domain
	dataDir := conf.Web.Subpath + conf.Web.AppDataPath
	content, err := dkim.MakeDkimConfFile(dataDir, domain)

	fmt.Println(content)
	fmt.Println(fmt.Sprintf("_dmarc in TXT ( v=DMARC1;p=quarantine;rua=mailto:admin@%s )", domain))
	fmt.Println(fmt.Sprintf("%s TXT ( v=spf1 a mx ~all )", domain))
	return err
}
