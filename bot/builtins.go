package bot

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"os"
	"strings"
	"time"

	otp "github.com/dgryski/dgoogauth"
	"github.com/ghodss/yaml"
)

// if help is more than tooLong lines long, send a private message
const tooLong = 14

// Size of QR code
const qrsize = 400

// If this list doesn't match what's registered below,
// you're gonna have a bad time.
var builtIns = []string{
	"builtInhelp",
	"builtInadmin",
	"builtIndump",
	"builtInlaunchcodes",
}

func init() {
	RegisterPlugin("builtIndump", PluginHandler{DefaultConfig: dumpConfig, Handler: dump})
	RegisterPlugin("builtInhelp", PluginHandler{DefaultConfig: helpConfig, Handler: help})
	RegisterPlugin("builtInadmin", PluginHandler{DefaultConfig: adminConfig, Handler: admin})
	RegisterPlugin("builtInlaunchcodes", PluginHandler{DefaultConfig: launchCodesConfig, Handler: launchCode})
}

/* builtin plugins, like help */

func launchCode(bot *Robot, command string, args ...string) {
	if command == "init" {
		return // ignore init
	}
	var userOTP otp.OTPConfig
	otpKey := "bot:OTP:" + bot.User
	updated := false
	lock, exists, ret := checkoutDatum(otpKey, &userOTP, true)
	if ret != Ok {
		bot.Say("Yikes - something went wrong with my brain, have somebody check my log")
		return
	}
	defer func() {
		if updated {
			ret := updateDatum(otpKey, lock, &userOTP)
			if ret != Ok {
				Log(Error, "Couldn't save OTP config")
				bot.Reply("Good grief. I'm having trouble remembering your launch codes - have somebody check my log")
			}
		} else {
			// Well-behaved plugins will always do a CheckinDatum when the datum hasn't been updated,
			// in case there's another thread waiting.
			checkinDatum(otpKey, lock)
		}
	}()
	switch command {
	case "send":
		if exists {
			bot.Reply("I've already sent you the launch codes, contact an administrator if you're having problems")
			return
		}
		otpb := make([]byte, 10)
		random.Read(otpb)
		userOTP.Secret = base32.StdEncoding.EncodeToString(otpb)
		userOTP.WindowSize = 2
		userOTP.DisallowReuse = []int{}
		var codeMail bytes.Buffer
		fmt.Fprintf(&codeMail, "For your authenticator:\n%s\n", userOTP.Secret)
		if ret := bot.Email("Your launch codes - if you print this email, please chew it up and swallow it", &codeMail); ret != Ok {
			bot.Reply("There was a problem sending your launch codes, contact an administrator")
			return
		}
		updated = true
		bot.Reply("I've emailed your launch codes - please delete it promptly")
	case "check":
		if !exists {
			bot.Reply("It doesn't appear you've been issued any launch codes, try 'send launch codes'")
			return
		}
		valid, err := userOTP.Authenticate(args[0])
		if err != nil {
			Log(Error, fmt.Errorf("Problem authenticating launch code: %v", err))
			bot.Reply("There was an error authenticating your launch code, have an adminstrator check the log")
			return
		}
		// Whether valid or not, the otp lib may update the struct
		updated = true
		if valid {
			bot.Reply("The launch code was valid")
		} else {
			bot.Reply("You supplied an invalid launch code, and I've contacted POTUS and the NSA")
		}
	}
}

func help(bot *Robot, command string, args ...string) {
	if command == "init" {
		return // ignore init
	}
	if command == "help" {
		b.lock.RLock()
		defer b.lock.RUnlock()

		var term, helpOutput string
		botSub := `(bot)`
		hasTerm := false
		helpLines := 0
		if len(args) == 1 && len(args[0]) > 0 {
			hasTerm = true
			term = args[0]
			if term == "help" {
				Log(Trace, "Help requested for help, returning")
				return
			}
			Log(Trace, "Help requested for term", term)
		}

		for _, plugin := range plugins {
			Log(Trace, fmt.Sprintf("Checking help for plugin %s (term: %s)", plugin.name, term))
			if !hasTerm { // if you ask for help without a term, you just get help for whatever commands are available to you
				if messageAppliesToPlugin(bot.User, bot.Channel, plugin) {
					for _, phelp := range plugin.Help {
						for _, helptext := range phelp.Helptext {
							if len(phelp.Keywords) > 0 && phelp.Keywords[0] == "*" {
								// * signifies help that should be prepended
								helpOutput = strings.Replace(helptext, botSub, b.name, -1) + string('\n') + helpOutput
							} else {
								helpOutput += strings.Replace(helptext, botSub, b.name, -1) + string('\n')
							}
							helpLines++
						}
					}
				}
			} else { // when there's a search term, give all help for that term, but add (channels: xxx) at the end
				for _, phelp := range plugin.Help {
					for _, keyword := range phelp.Keywords {
						if term == keyword {
							chantext := ""
							if plugin.DirectOnly {
								// Look: the right paren gets added below
								chantext = " (direct message only"
							} else {
								for _, pchan := range plugin.Channels {
									if len(chantext) == 0 {
										chantext += " (channels: " + pchan
									} else {
										chantext += ", " + pchan
									}
								}
							}
							if len(chantext) != 0 {
								chantext += ")"
							}
							for _, helptext := range phelp.Helptext {
								helpOutput += strings.Replace(helptext, botSub, b.name, -1) + chantext + string('\n')
								helpLines++
							}
						}
					}
				}
			}
		}
		if hasTerm {
			helpOutput = "Command(s) matching keyword: " + term + "\n" + helpOutput
		}
		switch {
		case helpLines == 0:
			bot.Say("Sorry, bub - I got nothin' for ya'")
		case helpLines > tooLong:
			if len(bot.Channel) > 0 {
				bot.Reply("(the help output was pretty long, so I sent you a private message)")
				if !hasTerm {
					helpOutput = "Command(s) available in channel: " + bot.Channel + "\n" + helpOutput
				}
			} else {
				if !hasTerm {
					helpOutput = "Command(s) available:" + "\n" + helpOutput
				}
			}
			bot.SendUserMessage(bot.User, strings.TrimRight(helpOutput, "\n"))
		default:
			if !hasTerm {
				helpOutput = "Command(s) available:" + "\n" + helpOutput
			}
			bot.Say(strings.TrimRight(helpOutput, "\n"))
		}
	}
}

func dump(bot *Robot, command string, args ...string) {
	if command == "init" {
		return // ignore init
	}
	b.lock.RLock()
	defer b.lock.RUnlock()
	switch command {
	case "robotdefault":
		bot.Fixed().Say(fmt.Sprintf("Here's my default configuration:\n%s", defaultConfig))
	case "robot":
		c, _ := yaml.Marshal(config)
		bot.Fixed().Say(fmt.Sprintf("Here's how I've been configured, irrespective of interactive changes:\n%s", c))
	case "plugdefault":
		if plug, ok := pluginHandlers[args[0]]; ok {
			bot.Fixed().Say(fmt.Sprintf("Here's the default configuration for \"%s\":\n%s", args[0], plug.DefaultConfig))
		} else { // look for an external plugin
			found := false
			for _, plugin := range plugins {
				if args[0] == plugin.name && plugin.pluginType == plugExternal {
					found = true
					if cfg, err := getExtDefCfg(plugin); err == nil {
						bot.Fixed().Say(fmt.Sprintf("Here's the default configuration for \"%s\":\n%s", args[0], *cfg))
					} else {
						bot.Say("I had a problem looking that up - somebody should check my logs")
					}
				}
			}
			if !found {
				bot.Say("Didn't find a plugin named " + args[0])
			}
		}
	case "plugin":
		found := false
		for _, plugin := range plugins {
			if args[0] == plugin.name {
				found = true
				c, _ := yaml.Marshal(plugin)
				bot.Fixed().Say(fmt.Sprintf("%s", c))
			}
		}
		if !found {
			bot.Say("Didn't find a plugin named " + args[0])
		}
	case "list":
		plist := make([]string, 0, len(plugins))
		for _, plugin := range plugins {
			plist = append(plist, plugin.name)
		}
		bot.Say(fmt.Sprintf("Here are the plugins I have configured:\n%s", strings.Join(plist, ", ")))
	}
}

var byebye = []string{
	"Sayonara!",
	"Adios",
	"Hasta la vista!",
	"Later gator!",
}

func admin(bot *Robot, command string, args ...string) {
	if command == "init" {
		return // ignore init
	}
	if !bot.CheckAdmin() {
		bot.Reply("Sorry, only an admin user can request that")
		return
	}
	switch command {
	case "reload":
		err := loadConfig()
		if err != nil {
			bot.Reply("Error encountered during reload, check the logs")
			Log(Error, fmt.Errorf("Reloading configuration, requested by %s: %v", bot.User, err))
			return
		}
		bot.Reply("Configuration reloaded successfully")
		Log(Info, "Configuration successfully reloaded by a request from:", bot.User)
	case "quit":
		// Get all important locks to make sure nothing is being changed - then exit
		botLock.Lock()
		b.lock.Lock()
		dataLock.Lock()
		bot.Say(bot.RandomString(byebye))
		Log(Info, "Exiting on administrator command")
		// How long does it _actually_ take for the message to go out?
		time.Sleep(time.Second)
		os.Exit(0)
	}
}
