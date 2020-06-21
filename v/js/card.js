var trial =	document.getElementById("trial"); 
var restart = document.getElementById("restart"); 
var hero = document.getElementById("hero");
var stage = document.getElementById("stage");
var speech = document.getElementById("speech");
var me = document.getElementById("me")

var hand0 = document.getElementById("hand0");
var hand1 = document.getElementById("hand1");
var hand2 = document.getElementById("hand2");
var hand3 = document.getElementById("hand3");
var hand4 = document.getElementById("hand4");
var hand5 = document.getElementById("hand5");
var hand6 = document.getElementById("hand6");
var hand7 = document.getElementById("hand7");

var hands = [hand0,hand1, hand2,hand3, hand4, hand5, hand6, hand7];

var Land0 = document.getElementById("Land0");
var Land1 = document.getElementById("Land1");
var Land2 = document.getElementById("Land2");
var Land3 = document.getElementById("Land3");
var Land4 = document.getElementById("Land4");
var Land5 = document.getElementById("Land5");
var Land6 = document.getElementById("Land6");

var Lands = [Land0, Land1, Land2, Land3, Land4, Land5, Land6];

var support0 = document.getElementById("support0");
var support1 = document.getElementById("support1");
var support2 = document.getElementById("support2");
var support3 = document.getElementById("support3");

var supports = [support0, support1, support2, support3 ]


var player0 = document.getElementById("player0");
var player1 = document.getElementById("player1");
var player2 = document.getElementById("player2");
var player3 = document.getElementById("player3");
var player4 = document.getElementById("player4");

var players = [player0 ,player1, player2, player3,player4]


var threat0 = document.getElementById("threat0");
var threat1 = document.getElementById("threat1");
var threat2 = document.getElementById("threat2");
var threat3 = document.getElementById("threat3");
var threat4 = document.getElementById("threat4");

var threats = [threat0, threat1, threat2, threat3, threat4]


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
				console.log("fasdf")
				Lands[i].style.backgroundImage = '';
			}
		}

		for (let i = 0; i < threats.length; i++) {
			var divs=threats[i].getElementsByTagName("div");
    		while(divs.length>0){
  		   	 threats[i].removeChild(divs[0]);
			}
		}

		for (let i = 0; i < jsonArray["Threats"].length; i++) {
			for (let j = 0; j < jsonArray["Threats"][i].length; j++) {
				var odiv=document.createElement("div");
				odiv.style = `background-image:url(./icon/${jsonArray["Threats"][i][j]}.png);`;
				threats[i].appendChild(odiv);
			}
		}

		for (let i = 0; i < jsonArray["Players"].length;i++){
			console.log(i)
			players[i].innerHTML = jsonArray["Players"][i]
		}
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

		for (let i = 0; i < jsonArray["Handcard"].length; i++) {
			hands[i].style = `background-image:url(./card/${jsonArray["Handcard"][i]["Name"]}.png);`;
		}
		hero.style = `background-image:url(./hero/${jsonArray["Hero"]["Name"]}.png)`;
		document.getElementById("speechNum").innerHTML="x"+jsonArray["SpeechTime"]
		support0.innerHTML = "x"+jsonArray["Support"]["Left"]
		support1.innerHTML = "x"+jsonArray["Support"]["Right"]
		support2.innerHTML = "x"+jsonArray["Support"]["Left2"]
		support3.innerHTML = "x"+jsonArray["Support"]["Right2"]
	}

	
};

trial.onclick = function () {
	if (stage.innerHTML!="Waiting") {
		let draw= '{"order":"draw"}';
		w.send(draw)
	}
}

function playCard(card){
	let playCard = '{"order":"playCard","choose":"'+card+'"}';
	w.send(playCard);
}

hand0.onclick=function(){
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("0")
	}
}
hand1.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("1")
	}
};
hand2.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("2");
	}
};
hand3.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("3")
	}
};
hand4.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("4")
	}
};
hand5.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("5")
	}
};
hand6.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("6")
	}
};
hand7.onclick = function () {
	if (this.style.backgroundImage.search("cardBack")==-1) {
		playCard("7")
	}
};

hero.onclick = function () {
	if (this.style.backgroundImage.search("used") == -1) {
		let heroUse= '{"order":"heroUse"}';
		w.send(heroUse);
	}
}

function lucky(land){
	if (stage.innerHTML=="幸運草") {
		let lucky = '{"order":"luckyClover","choose":"'+land+'"}';
		w.send(lucky)
	}
}

Land0.onclick = function () {
	lucky(0)
};
Land1.onclick = function () {
	lucky(1)
};
Land2.onclick = function () {
	lucky(2)
};
Land3.onclick = function () {
	lucky(3)
};
Land4.onclick = function () {
	lucky(4)
};
Land5.onclick = function () {
	lucky(5)
};
Land6.onclick = function () {
	lucky(6)
};

restart.onclick = function () {
	let restart= '{"order":"restart"}';
	w.send(restart);
}

speech.onclick = function () {

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


