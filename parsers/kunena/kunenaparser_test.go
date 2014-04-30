package kunena

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var html = `<td class="kmessage-left">
				
<div class="kmsgbody">
	<div class="kmsgtext">
		<i>Vietnam, 18. April 1965<br />
07:30 Uhr Ortszeit</i><br />
<br />
Charlie hat sich in unserem Einsatzgebiet erheblich verbreitet. Hinterhältige Angriffe auf unsere Basis und unsere Truppen sind unser ständiger Begleiter.<br />
Wir werden in unserem Einsatzgebiet eine Flusspatrouille durchführen, um Charlie in diesem Gebiet endgültig auszuräuchern.<br />
<br />
Sie starten im <i>Camp Oscar</i>. Von dortaus führen Sie die <i>Patrouille </i>in Richtung <i>Westen </i>Entlang des Flusses. Ihr Zielpunkt wir der Ort <i>Phang Ra</i> sein, wo wir eine Basis des Feindes vermuten.<br />
<br />
Des weiteren vermissen wir einen Huey, der auf einem Versorgungsflug war. Wir vermuten, dass dieser abgestürzt ist. Finden Sie den Helikopter auf Ihrem Weg.<br />
<br />
Meine Herren, wir sind nicht mehr hier, um uns in diesem Gebiet Freunde zu machen! Dieser Krieg dauert einfach schon zu lange und hat schon zu viele Verwundete und Tote gekostet. Zeigen Sie Charlie, dass wir die überlegenere Armee sind!<br />
<br />
Einsatzunterstützung:<br />
<br />
- UH- 1F Gunship zur CAS- Unterstützung<br />
<br />
<br />
<b>03.05.2014<br />
19:30 Treffen, Gruppenführerbesprechung<br />
20:00 Kampfeinsatz<br />
<br />
Modstring:</b><br />
@e12_dac;@e12_tools;@he_acex_usnavy;@he_acex_ru;@he_acex;@he_acre;@he_beta;@he_beta\expansion;@he_cba_co;@he_jayarma2lib;@he_ace;@he_blastcore;@he_unsung;<br />
<br />
<b>Slots: </b><br />
<br />
#1 Platoonleader: <span style="color:#ffbb00">HG2012 Trigger</span><br />
<br />
<b>1st Heavy Fireteam (Reserviert für Heeresgruppe 2012)</b><br />
<br />
#2 Fireteamleader: <span style="color:#ffbb00">HG2012 Spawnferkel</span><br />
#3 Designated Marksman: <span style="color:#ffbb00">HG2012 Reaper</span><br />
#4 Soldier: <span style="color:#ffbb00">HG2012 00Zyan00</span><br />
#5 Corpsman: <span style="color:#ffbb00">HG2012 JohnBerg</span><br />
#6 Automatic Rifleman: <span style="color:#ffbb00"></span><br />
#7 Soldier: <span style="color:#ffbb00">HG2012 Smula</span><br />
#8 Pionier: <span style="color:#ffbb00">HG2012 Möwe</span><br />
#9 LAW- Soldier: <span style="color:#ffbb00">HG2012 Dragonmaster</span><br />
#Zusatz1 Automatic Rifleman: <span style="color:#ffbb00"></span><br />
#Zusatz2 Soldier: <span style="color:#ffbb00"></span><br />
<br />
<b>2nd Heavy Fireteam (Reserviert für Echo 12)</b><br />
<br />
#10 Fireteamleader: <span style="color:#ffbb00">E12 Kerodan</span><br />
#11 Designated Marksman: <span style="color:#ffbb00">E12 Badbug</span><br />
#12 Soldier: <span style="color:#ffbb00">E12 RIPchen</span><br />
#13 Corpsman: <span style="color:#ffbb00">E12 Devilous</span><br />
#14 Automatic Rifleman: <span style="color:#ffbb00">E12 Rickyfox</span><br />
#15 Soldier: <span style="color:#ffbb00">E12 Proofy</span><br />
#16 Pionier: <span style="color:#ffbb00">E12 Systemstörung81</span><br />
#17 LAW- Soldier: <span style="color:#ffbb00">E12 Spirit</span><br />
#Zusatz3 Automatic Rifleman: <span style="color:#ffbb00"></span><br />
#Zusatz4 Soldier: <span style="color:#ffbb00"></span><br />
<br />
<b>UH- 1F Gunship (In Absprache mit mir!)</b><br />
<br />
#18 Pilot: <span style="color:#ffbb00">E12 Bob</span><br />
#19 Pilot: <span style="color:#ffbb00">E12 Pogo</span><br />
#20 Doorgunner: <span style="color:#ffbb00">HG2012 Atze</span><br />
#21 Doorgunner: <span style="color:#ffbb00">HG2012 Neville</span>	</div>
</div>
</div>
</td>

`

func TestParseSlotlist(t *testing.T) {
	Convey("Given a fresh parser", t, func() {

		p := &KunenaParser{}
		Convey("Parser accepts url", func() {
			accept := p.Accept("http://heeresgruppe2012.de/index.php/forum/events-und-missionen/252-03-05-2014-co-20-operation-baker")
			So(accept, ShouldBeTrue)
		})
		Convey("When parsing wiki slotlist", func() {
			sl := p.Parse(html, "http://heeresgruppe2012.de/index.php/forum/events-und-missionen/252-03-05-2014-co-20-operation-baker")
			Convey("Slotlist can be parsed", func() {
				So(sl, ShouldNotBeNil)
			})

			Convey("Has 4 groups", func() {
				So(len(sl.SlotListGroups), ShouldEqual, 4)
				Convey("Second Group", func() {
					g := sl.SlotListGroups[1]
					Convey("Has 10 slots", func() {
						So(len(g.Slots), ShouldEqual, 10)
					})
				})
			})

		})
	})
}
