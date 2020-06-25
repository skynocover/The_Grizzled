var trial =	document.getElementById("trial"); 
var morale =document.getElementById("morale"); 
var restart = document.getElementById("restart"); 
var hero = document.getElementById("hero");
var stage = document.getElementById("stage");
var speech = document.getElementById("speech");
var me = document.getElementById("me");

function getElementsByIdStartsWith(container, selectorTag, prefix) {
    var items = [];
    var myPosts = document.getElementById(container).getElementsByTagName(selectorTag);
    for (var i = 0; i < myPosts.length; i++) {
        if (myPosts[i].id.lastIndexOf(prefix, 0) === 0) {
            items.push(myPosts[i]);
        }
	}
    return items;
}

hands = getElementsByIdStartsWith("hands", "div", "hand")
for (let i = 0; i < hands.length; i++) {
	hands[i].onclick= function () {
		if (hands[i].style.backgroundImage.search("cardBack")==-1) {
			if (stage.innerHTML.split(":")[0]=="演說" ){
				let playCard = '{"order":"speechCard","choose":"'+i+'"}';
				w.send(playCard);
			}else{
				let playCard = '{"order":"playCard","choose":"'+i+'"}';
				w.send(playCard);
			}
		}
	}
}

Lands = getElementsByIdStartsWith("lands", "div", "Land")
for (let i = 0; i < Lands.length; i++) {
	Lands[i].onclick= function () {
		if (stage.innerHTML=="幸運草") {
			let lucky = '{"order":"luckyClover","choose":"'+i+'"}';
			w.send(lucky)
		}
	}
}

supports = getElementsByIdStartsWith("supports", "div", "support")
for (let i = 0; i < supports.length; i++) {
	supports[i].onclick= function () {
		let sup = '{"order":"support","choose":"'+i+'"}';
		w.send(sup)
	}
}

players = getElementsByIdStartsWith("players", "span", "player")
threats = getElementsByIdStartsWith("players", "div", "threat")

w = new WebSocket("ws://" + HOST + "/my_endpoint");
/*
w.onopen = function () {
	console.log("Websocket connection enstablished");
};

w.onclose = function () {
	appendMessage("<div><center><h3>Disconnected</h3></center></div>");
};
*/
w.onmessage = function (message) {
	var jsonArray=JSON.parse(message.data);
	console.log(jsonArray)

	if (jsonArray["Process"]=="game") {

		for (let i = 0; i < 7; i++) {
			if (i < jsonArray["NoMansLand"].length) {
				Lands[i].style = `background-image:url(./card/${jsonArray["NoMansLand"][i]["Name"]}.png);`;
			}else{
				Lands[i].style.backgroundImage = '';
			}
		}
		// 清空threat
		for (let i = 0; i < threats.length; i++) {
			var divs=threats[i].getElementsByTagName("div");
    		while(divs.length>0){
  		   	 threats[i].removeChild(divs[0]);
			}
		}
		//放入threat
		for (let i = 0; i < jsonArray["Threats"].length; i++) {
			for (let j = 0; j < jsonArray["Threats"][i].length; j++) {
				var odiv=document.createElement("div");
				odiv.style = `background-image:url(./icon/${jsonArray["Threats"][i][j]}.png);`;
				threats[i].appendChild(odiv);
			}
		}

		for (let i = 0; i < jsonArray["Players"].length;i++){
			players[i].innerHTML = jsonArray["Players"][i]
		}

		trialnum = jsonArray["TM"][0]
		moralenum = jsonArray["TM"][1]

		stage.innerHTML = jsonArray["Stage"]
		if (jsonArray["Stage"]=="Start!") {
			swal({
  				title: "遊戲就緒",
  				icon: "success",
  				button: "抽牌",
			})
			.then(() =>{
				let draw= '{"order":"draw"}';
				w.send(draw)
			});
		}
	}

	if (jsonArray["Process"]=="hand") {
		
		for (let i = 0; i<8; i++) {
			if (i < jsonArray["Handcard"].length) {
				hands[i].style = `background-image:url(./card/${jsonArray["Handcard"][i]["Name"]}.png);`;
			}else{
				hands[i].style = `background-image:url(./card/cardBack.png);`;
			}
		}
		hero.style = `background-image:url(./hero/${jsonArray["Hero"]["Name"]}.png)`;
		
		document.getElementById("speechNum").innerHTML="x"+jsonArray["SpeechTime"]
		support0.innerHTML = "x"+jsonArray["Support"]["Left"]
		support1.innerHTML = "x"+jsonArray["Support"]["Right"]
		support2.innerHTML = "x"+jsonArray["Support"]["Left2"]
		support3.innerHTML = "x"+jsonArray["Support"]["Right2"]
	}
};


hero.onclick = function () {
	if (this.style.backgroundImage.search("used") == -1) {
		let heroUse= '{"order":"heroUse"}';
		w.send(heroUse);
	}
}

restart.onclick = function () {
	let restart= '{"order":"restart"}';
	w.send(restart);
}

speech.onclick = function () {
	if (this.innerHTML!="x0") {
		swal("選擇一個威脅", {
  		buttons: {
  			cancel: "Cancel!",
    		take0:{text:"雨天",value:"Rain"},
    		take1:{text:"雪天",value:"Snow"},
    		take2:{text:"夜晚",value:"Night"},
    		take3:{text:"子彈",value:"Bullet"},
    		take4:{text:"哨子",value:"Whistle"},
    		take5:{text:"面具",value:"Mask"},
  		},
		}).then((value) => {
		if (value!=null) {
			let speech = '{"order":"speech","choose":"'+value+'"}';
			w.send(speech);
		}
		})
	}
}





/* 牌庫檢查 */
var trialnum = 0
trial.onclick = function () {
	swal({
  		title: "考驗:"+trialnum,
	});
}

var moralenum = 0
morale.onclick = function () {
	swal({
  		title: "士氣:"+moralenum,
	});
}
/* 牌庫檢查結束 */
