package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMultiREOutput(t *testing.T) {
	body := `<rpc-reply xmlns:junos="http://xml.juniper.net/junos/17.2R1/junos">
			<multi-routing-engine-results>
				<multi-routing-engine-item>
					<re-name>fpc0</re-name>
					<system-storage-information junos:style="brief">
						<filesystem>
							<filesystem-name>/dev/gpt/junos</filesystem-name>
							<total-blocks junos:format="1.3G">2796512</total-blocks>
							<used-blocks junos:format="814M">1667792</used-blocks>
							<available-blocks junos:format="442M">905000</available-blocks>
							<used-percent> 65</used-percent>
							<mounted-on>/.mount</mounted-on>
						</filesystem>
						<filesystem>
							<filesystem-name>/dev/sda</filesystem-name>
							<total-blocks junos:format="1.3G">2796512</total-blocks>
							<used-blocks junos:format="814M">1667792</used-blocks>
							<available-blocks junos:format="442M">905000</available-blocks>
							<used-percent> 6r75</used-percent>
							<mounted-on>/</mounted-on>
						</filesystem>
					</system-storage-information>
				</multi-routing-engine-item>
				<multi-routing-engine-item>
					<re-name>fpc1</re-name>
					<system-storage-information junos:style="brief">
						<filesystem>
							<filesystem-name>/dev/gpt/junos1</filesystem-name>
							<total-blocks junos:format="1.1G">2796512</total-blocks>
							<used-blocks junos:format="810M">1667792</used-blocks>
							<available-blocks junos:format="440M">905000</available-blocks>
							<used-percent> 65</used-percent>
							<mounted-on>/.mount</mounted-on>
						</filesystem>
					</system-storage-information>
				</multi-routing-engine-item>

			</multi-routing-engine-results>
			<cli>
				<banner>{master:0}</banner>
			</cli>
		</rpc-reply>`

	rpc := RpcReply{}
	err := parseXML([]byte(body), &rpc)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, rpc.MultiRoutingEngineResults.RoutingEngine[0].StorageInformation)

	// test first routing engine
	assert.Equal(t, "fpc0", rpc.MultiRoutingEngineResults.RoutingEngine[0].Name, "re-name")

	f := rpc.MultiRoutingEngineResults.RoutingEngine[0].StorageInformation.Filesystems[1]

	assert.Equal(t, "/dev/sda", f.FilesystemName, "filesystem-name")
	assert.Equal(t, int64(2796512), f.TotalBlocks, "total-blocks")
	assert.Equal(t, int64(1667792), f.UsedBlocks, "used-blocks")
	assert.Equal(t, "/", f.MountedOn, "mounted-on")

	// test second routing engine
	assert.Equal(t, "fpc1", rpc.MultiRoutingEngineResults.RoutingEngine[1].Name, "re-name")

	f = rpc.MultiRoutingEngineResults.RoutingEngine[1].StorageInformation.Filesystems[0]

	assert.Equal(t, "/dev/gpt/junos1", f.FilesystemName, "filesystem-name")
	assert.Equal(t, int64(2796512), f.TotalBlocks, "total-blocks")
	assert.Equal(t, int64(1667792), f.UsedBlocks, "used-blocks")
	assert.Equal(t, "/.mount", f.MountedOn, "mounted-on")
}
