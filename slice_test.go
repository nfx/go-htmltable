package htmltable

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type nice struct {
	C string `header:"c"`
	D string `header:"d"`
}

func TestNewSliceFromString(t *testing.T) {
	out, err := NewSliceFromString[nice](fixture)
	assertNoError(t, err)
	assertEqual(t, []nice{
		{"2", "5"},
		{"4", "6"},
	}, out)
}

type Ticker struct {
	Symbol   string `header:"Symbol"`
	Security string `header:"Security"`
	CIK      string `header:"CIK"`
}

func TestNewSliceFromUrl(t *testing.T) {
	url := "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"
	out, err := NewSliceFromURL[Ticker](url)
	assertNoError(t, err)
	assertGreaterOrEqual(t, len(out), 500)
}

func TestNewSliceFromUrl_Fails(t *testing.T) {
	_, err := NewSliceFromURL[Ticker]("https://127.0.0.1")
	assertEqualError(t, err, "Get \"https://127.0.0.1\": dial tcp 127.0.0.1:443: connect: connection refused")
}

func TestNewSliceFromUrl_NoTables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()
	_, err := NewSliceFromURL[Ticker](server.URL)
	assertEqualError(t, err, "cannot find table with columns: Symbol, Security, CIK")
}

func TestNewSliceInvalidTypes(t *testing.T) {
	type exotic struct {
		A string  `header:""`
		C float32 `header:"c"`
	}
	_, err := NewSliceFromString[exotic](fixture)
	assertEqualError(t, err, "only strings are supported, C is float32")
}

func TestIssue3(t *testing.T) {
	type Chipset struct {
		Model       string `header:"Model"`
		Codename    string `header:"Codename"`
		Released    string `header:"Released"`
		CPUSupport  string `header:"CPU support"`
		Fab         string `header:"Fab (nm)"`
		HT          string `header:"HT (MHz)"`
		IGP         string `header:"IGP"`
		CrossFire   string `header:"CrossFire"`
		Southbridge string `header:"Southbridge"`
		Features    string `header:"Features / Notes"`
	}
	chipsets, err := NewSliceFromString[Chipset](x)
	assertNoError(t, err)
	assertEqual(t, true, len(chipsets) > 3)
}

const x = `<table class="wikitable" style="text-align:center">
<tbody>
   <tr>
	  <th rowspan="2">Model</th>
	  <th rowspan="2">Release date</th>
	  <th rowspan="2"><a href="/wiki/PCI_Express" title="PCI Express">PCIe</a> support<sup id="cite_ref-pcie_20-0" class="reference"><a href="#cite_note-pcie-20">[a]</a></sup></th>
	  <th colspan="2">Multi-GPU</th>
	  <th rowspan="2"><a href="/wiki/USB" title="USB">USB</a> support<sup id="cite_ref-usb_21-0" class="reference"><a href="#cite_note-usb-21">[b]</a></sup></th>
	  <th colspan="3">Storage features</th>
	  <th rowspan="2">Processor<br><a href="/wiki/Overclocking" title="Overclocking">overclocking</a></th>
	  <th rowspan="2"><a href="/wiki/Thermal_design_power" title="Thermal design power">TDP</a></th>
	  <th colspan="5">CPU support<sup id="cite_ref-22" class="reference"><a href="#cite_note-22">[14]</a></sup></th>
	  <th rowspan="2">Architecture</th>
   </tr>
   <tr>
	  <th><a href="/wiki/AMD_CrossFireX" class="mw-redirect" title="AMD CrossFireX">CrossFire</a></th>
	  <th><a href="/wiki/Scalable_Link_Interface" title="Scalable Link Interface">SLI</a></th>
	  <th><a href="/wiki/Serial_ATA" title="Serial ATA">SATA</a> ports</th>
	  <th><a href="/wiki/RAID" title="RAID">RAID</a></th>
	  <th><a rel="nofollow" class="external text" href="https://www.amd.com/en/technologies/store-mi">AMD StoreMI</a></th>
	  <th><a href="/wiki/Excavator_(microarchitecture)" title="Excavator (microarchitecture)">Excavator</a></th>
	  <th><a href="/wiki/Zen_(first_generation_microarchitecture)" class="mw-redirect" title="Zen (first generation microarchitecture)">Zen</a></th>
	  <th><a href="/wiki/Zen%2B" title="Zen+">Zen+</a></th>
	  <th><a href="/wiki/Zen_2" title="Zen 2">Zen 2</a></th>
	  <th><a href="/wiki/Zen_3" title="Zen 3">Zen 3</a></th>
   </tr>
   <tr>
	  <th>A320</th>
	  <td>February 2017<sup id="cite_ref-:date3s_23-0" class="reference"><a href="#cite_note-:date3s-23">[15]</a></sup></td>
	  <td>PCIe 2.0 ×4</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td rowspan="2" style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td>1, 2, 6</td>
	  <td rowspan="2">4</td>
	  <td rowspan="8">0,<br>1,<br>10</td>
	  <td rowspan="3" style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td style="background:#FFB;vertical-align:middle;text-align:center;" class="table-partial">Limited to pre-Zen CPUs, unless an unsupported third-party motherboard firmware applied</td>
	  <td rowspan="7">~5&nbsp;W<sup id="cite_ref-24" class="reference"><a href="#cite_note-24">[16]</a></sup></td>
	  <td rowspan="3" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="5" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="5" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="3" colspan="2" style="background: #FF8; vertical-align: middle; text-align: center;" class="table-maybe"><small>Varies</small><sup id="cite_ref-zen2_25-0" class="reference"><a href="#cite_note-zen2-25">[c]</a></sup></td>
	  <td rowspan="7">Promontory</td>
   </tr>
   <tr>
	  <th>B350</th>
	  <td>February 2017<sup id="cite_ref-:date3s_23-1" class="reference"><a href="#cite_note-:date3s-23">[15]</a></sup></td>
	  <td>PCIe 2.0 ×6</td>
	  <td rowspan="4" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td>2, 2, 6</td>
	  <td rowspan="2" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
   </tr>
   <tr>
	  <th>X370</th>
	  <td>February 2017<sup id="cite_ref-:date3s_23-2" class="reference"><a href="#cite_note-:date3s-23">[15]</a></sup></td>
	  <td>PCIe 2.0 ×8</td>
	  <td style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td>2, 6, 6</td>
	  <td>8</td>
   </tr>
   <tr>
	  <th>B450</th>
	  <td>March 2018<sup id="cite_ref-:date4s_26-0" class="reference"><a href="#cite_note-:date4s-26">[17]</a></sup></td>
	  <td>PCIe 2.0 ×6</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td>2, 2, 6</td>
	  <td>4</td>
	  <td rowspan="5" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="2" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes,<br>with <abbr title="Precision Boost Overdrive">PBO</abbr></td>
	  <td rowspan="4" style="background: #FF8; vertical-align: middle; text-align: center;" class="table-maybe"><small>Varies</small><sup id="cite_ref-zen3_27-0" class="reference"><a href="#cite_note-zen3-27">[d]</a></sup></td>
	  <td rowspan="5" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="4" style="background: #FF8; vertical-align: middle; text-align: center;" class="table-maybe"><small>Varies</small><sup id="cite_ref-zen3_27-1" class="reference"><a href="#cite_note-zen3-27">[d]</a></sup><sup id="cite_ref-Tom's_Hardware_Zen3_Update_28-0" class="reference"><a href="#cite_note-Tom's_Hardware_Zen3_Update-28">[18]</a></sup></td>
   </tr>
   <tr>
	  <th>X470</th>
	  <td>March 2018<sup id="cite_ref-:date4s_26-1" class="reference"><a href="#cite_note-:date4s-26">[17]</a></sup></td>
	  <td>PCIe 2.0 ×8</td>
	  <td style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td>2, 6, 6</td>
	  <td>8</td>
   </tr>
   <tr>
	  <th>A520</th>
	  <td>August 2020<sup id="cite_ref-29" class="reference"><a href="#cite_note-29">[19]</a></sup></td>
	  <td>PCIe 3.0 ×6</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td>1, 2, 6</td>
	  <td>4</td>
	  <td style="background:#FFB;vertical-align:middle;text-align:center;" class="table-partial">No, unless an unsupported third-party motherboard firmware applied</td>
	  <td rowspan="2" style="background: #FF8; vertical-align: middle; text-align: center;" class="table-maybe">Varies</td>
	  <td rowspan="3" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
   </tr>
   <tr>
	  <th>B550<sup id="cite_ref-30" class="reference"><a href="#cite_note-30">[e]</a></sup></th>
	  <td>June 2020<sup id="cite_ref-31" class="reference"><a href="#cite_note-31">[20]</a></sup></td>
	  <td>PCIe 3.0 ×10<sup id="cite_ref-32" class="reference"><a href="#cite_note-32">[21]</a></sup></td>
	  <td rowspan="2" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td style="background: #FF8; vertical-align: middle; text-align: center;" class="table-maybe">Varies</td>
	  <td>2, 2, 6</td>
	  <td>6</td>
	  <td rowspan="2" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes,<br>with <abbr title="Precision Boost Overdrive">PBO</abbr></td>
   </tr>
   <tr>
	  <th>X570</th>
	  <td>July 2019<sup id="cite_ref-33" class="reference"><a href="#cite_note-33">[22]</a></sup></td>
	  <td>PCIe 4.0 ×16</td>
	  <td style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td>8, 0, 4</td>
	  <td>12</td>
	  <td>~15&nbsp;W<sup id="cite_ref-34" class="reference"><a href="#cite_note-34">[23]</a></sup><sup id="cite_ref-35" class="reference"><a href="#cite_note-35">[24]</a></sup> <sup id="cite_ref-36" class="reference"><a href="#cite_note-36">[f]</a></sup></td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No <sup id="cite_ref-x570_37-0" class="reference"><a href="#cite_note-x570-37">[g]</a></sup></td>
	  <td style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td>Bixby</td>
   </tr>
</tbody>
</table>`

const unparsable = `<table class="wikitable" style="font-size: 100%; text-align: center; letter-spacing:0px;">
<tbody>
   <tr>
	  <th rowspan="2">Branding</th>
	  <th rowspan="2">Release date</th>
	  <th rowspan="2">Chipsets</th>
	  <th colspan="2">Chipset Links</th>
	  <th rowspan="2"><a href="/wiki/PCI_Express" title="PCI Express">PCIe</a> support<sup id="cite_ref-pcie_48-0" class="reference"><a href="#cite_note-pcie-48">[a]</a></sup></th>
	  <th colspan="2">Multi-GPU</th>
	  <th colspan="3"><a href="/wiki/USB" title="USB">USB</a> support</th>
	  <th colspan="2">Storage features</th>
	  <th rowspan="2">Processor<br><a href="/wiki/Overclocking" title="Overclocking">overclocking</a></th>
	  <th rowspan="2"><a href="/wiki/Thermal_design_power" title="Thermal design power">TDP</a></th>
	  <th colspan="1">CPU support</th>
   </tr>
   <tr>
	  <th>CPU</th>
	  <th>Interchipset</th>
	  <th><a href="/wiki/AMD_CrossFireX" class="mw-redirect" title="AMD CrossFireX">CrossFire</a></th>
	  <th><a href="/wiki/Scalable_Link_Interface" title="Scalable Link Interface">SLI</a></th>
	  <th><a href="/wiki/USB#2.0" title="USB">2.0</a></th>
	  <th><a href="/wiki/USB#3.x" title="USB">3.2 Gen 2</a></th>
	  <th>Additional</th>
	  <th><a href="/wiki/RAID" title="RAID">RAID</a></th>
	  <th><a href="/wiki/Serial_ATA" title="Serial ATA">SATA III</a></th>
	  <th><a href="/wiki/Zen_4" title="Zen 4">Zen 4</a></th>
   </tr>
   <tr>
	  <th>A620</th>
	  <td>2023</td>
	  <td>Promontory 21<br>×1</td>
	  <td rowspan="5">PCIe 4.0 ×4</td>
	  <td rowspan="3" style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">Unused</td>
	  <td style="background: #EEE; font-size: smaller; vertical-align: middle; text-align: center;" class="unknown table-unknown">Un­known</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td rowspan="5" style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td colspan="5" style="background: #EEE; font-size: smaller; vertical-align: middle; text-align: center;" class="unknown table-unknown">Un­known</td>
	  <td style="background:#FFC7C7;vertical-align:middle;text-align:center;" class="table-no">No</td>
	  <td>~4.5W</td>
	  <td rowspan="5" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
   </tr>
   <tr>
	  <th>B650</th>
	  <td rowspan="2">October 10, 2022</td>
	  <td rowspan="2">Promontory 21<br>×1</td>
	  <td rowspan="2">PCIe 4.0 ×8<br>PCIe 3.0 ×4<sup id="cite_ref-49" class="reference"><a href="#cite_note-49">[35]</a></sup></td>
	  <td rowspan="4" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="2">×6<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-0" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2">×4<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-1" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2">×1 3.2 Gen 2×2<br><i>OR</i><br>×2 3.2 Gen 2<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-2" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2"></td>
	  <td rowspan="2">4</td>
	  <td rowspan="4" style="background:#9EFF9E;vertical-align:middle;text-align:center;" class="table-yes">Yes</td>
	  <td rowspan="2">~7W</td>
   </tr>
   <tr>
	  <th>B650E</th>
   </tr>
   <tr>
	  <th>X670</th>
	  <td rowspan="2">September 27, 2022</td>
	  <td rowspan="2">Promontory 21<br>×2</td>
	  <td rowspan="2">PCIe 4.0 ×4</td>
	  <td rowspan="2">PCIe 4.0 ×12<br>PCIe 3.0 ×8</td>
	  <td rowspan="2">×12<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-3" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2">×8<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-4" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2">×2 3.2 Gen 2×2<br><i>OR</i><br> ×1 3.2 Gen 2×2 <br>+<br> ×2 3.2 Gen 2<br><i>OR</i><br> ×4 3.2 Gen 2<sup id="cite_ref-ArsTechnica_Zen4Chipset_50-5" class="reference"><a href="#cite_note-ArsTechnica_Zen4Chipset-50">[36]</a></sup></td>
	  <td rowspan="2"></td>
	  <td rowspan="2">8<sup id="cite_ref-51" class="reference"><a href="#cite_note-51">[b]</a></sup></td>
	  <td rowspan="2">~14W<sup id="cite_ref-52" class="reference"><a href="#cite_note-52">[c]</a></sup></td>
   </tr>
   <tr>
	  <th>X670E</th>
   </tr>
</tbody>
</table>`