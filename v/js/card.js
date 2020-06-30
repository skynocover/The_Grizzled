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
			}else if(me.innerHTML==stage.innerHTML){
				let playCard = '{"order":"playCard","choose":"'+i+'"}';
				w.send(playCard);
			}
		}
	}
}

Lands = getElementsByIdStartsWith("lands", "div", "Land")
for (let i = 0; i < Lands.length; i++) {
	Lands[i].onclick= function () {
		if (stage.innerHTML.split(":")[0]==me.innerHTML &&  stage.innerHTML.split(":")[1]=="幸運草") {
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
				w.send('{"order":"draw"}')
			});
		}else if (jsonArray["Stage"]=="SupportSkip"){
			if (jsonArray["RoundWin"]) {
				mission = "任務成功"
				missionicon = "success"
			}else{
				mission = "任務失敗"
				missionicon = "error"
			}


			swal({
  				title: mission,
  				icon: missionicon,
			}).then(() =>{
				swal({
  					title: "沒有人得到支援",
  					icon: "error",
  					button: "確認",
				})
				.then((value) =>{
					w.send('{"order":"supportEnd","choose":""}');
				});
			});

		}else if (jsonArray["Stage"].split(":")[0]=="Support"){ //支援
			if (jsonArray["RoundWin"]) {
				mission = "任務成功"
				missionicon = "success"
			}else{
				mission = "任務失敗"
				missionicon = "error"
			}

			if ((jsonArray["Stage"].split(":")[1]==me.innerHTML)) {
				swal({
  					title: mission,
  					icon: missionicon,
				}).then(() =>{
					swal({
  						title: "得到支援",
  						icon: "success",
  						buttons: {
    						take0:{text:"消除身上威脅",value:"Threat"},
    						take1:{text:"回覆英雄能力",value:"Lucky"},
  						},
					})
					.then((value) =>{
						w.send('{"order":"supportEnd","choose":"'+value+'"}');
					});
				});
				
			}else{
				swal({
  					title: mission,
  					icon: missionicon,
				}).then(() =>{
					swal({
  						title: jsonArray["Stage"].split(":")[1]+"得到支援",
  						icon: "info",
  						button: "確認",
					})
					.then(() =>{
						w.send('{"order":"supportEnd","choose":""}')
					});
				});
			}

			/*
			if ((jsonArray["Stage"].split(":")[1]==me.innerHTML)) {
				swal({
  					title: "得到支援",
  					icon: "success",
  					buttons: {
    					take0:{text:"消除身上威脅",value:"Threat"},
    					take1:{text:"回覆英雄能力",value:"Lucky"},
  					},
				})
				.then((value) =>{
					w.send('{"order":"supportEnd","choose":"'+value+'"}');
				});
			}else{
				swal({
  					title: jsonArray["Stage"].split(":")[1]+"得到支援",
  					icon: "info",
  					button: "確認",
				})
				.then(() =>{
					w.send('{"order":"supportEnd","choose":""}')
				});
			}
			*/
			
		}else if  (jsonArray["Stage"].split(":")[0]=="Leader" && jsonArray["Stage"].split(":")[1]==me.innerHTML){
			swal("選擇抽牌張數",{
  				buttons: {
    				take0:{text:"抽一張",value:"1"},
    				take1:{text:"抽兩張",value:"2"},
    				take2:{text:"抽三張",value:"3"},
    				take3:{text:"抽四張",value:"4"},
  				},
			}).then((value) => {
				w.send('{"order":"newRound","choose":"'+value+'"}')
			})
		}else if (jsonArray["Stage"]=="WinGame"){
			swal("遊戲結束!", "和平再現!", "success");
		}else if (jsonArray["Stage"]=="LoseGame"){
			swal("遊戲結束!", "紀念碑", "error");
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
		support0.innerHTML = "x"+jsonArray["Support"][0]
		support1.innerHTML = "x"+jsonArray["Support"][1]
		support2.innerHTML = "x"+jsonArray["Support"][2]
		support3.innerHTML = "x"+jsonArray["Support"][3]
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
