var stage = document.getElementById('stage');
var me = document.getElementById('me');

function getElementsByIdStartsWith(container, selectorTag, prefix) {
  var items = [];
  var myPosts = document
    .getElementById(container)
    .getElementsByTagName(selectorTag);
  for (var i = 0; i < myPosts.length; i++) {
    if (myPosts[i].id.lastIndexOf(prefix, 0) === 0) {
      items.push(myPosts[i]);
    }
  }
  return items;
}

/////////////////////////////////////////////////////////////

hands = getElementsByIdStartsWith('hands', 'div', 'hand');
for (let i = 0; i < hands.length; i++) {
  hands[i].onclick = function () {
    if (hands[i].style.backgroundImage.search('cardBack') == -1) {
      if (stage.innerHTML.split(':')[0] == '演說') {
        let playCard = '{"order":"speechCard","choose":"' + i + '"}';
        w.send(playCard);
      } else if (me.innerHTML == stage.innerHTML) {
        let playCard = '{"order":"playCard","choose":"' + i + '"}';
        w.send(playCard);
      }
    }
  };
}

lands = getElementsByIdStartsWith('lands', 'div', 'Land');
for (let i = 0; i < lands.length; i++) {
  lands[i].onclick = function () {
    if (
      stage.innerHTML.split(':')[0] == me.innerHTML &&
      stage.innerHTML.split(':')[1] == '幸運草'
    ) {
      let lucky = '{"order":"luckyClover","choose":"' + i + '"}';
      w.send(lucky);
    }
  };
}

supports = getElementsByIdStartsWith('supports', 'div', 'support');
for (let i = 0; i < supports.length; i++) {
  supports[i].onclick = function () {
    if (this.innerHTML != 'x0' && me.innerHTML == stage.innerHTML) {
      let sup = '{"order":"support","choose":"' + i + '"}';
      w.send(sup);
    }
  };
}

var hero = document.getElementById('hero');
hero.onclick = function () {
  if (
    this.style.backgroundImage.search('used') == -1 &&
    me.innerHTML == stage.innerHTML
  ) {
    swal({
      title: '使用英雄能力?',
      text: '可從場上移除一張符合英雄能力的威脅',
      icon: 'warning',
      buttons: true,
      dangerMode: true,
    }).then((heroUse) => {
      if (heroUse) {
        w.send('{"order":"heroUse"}');
      }
    });
  }
};

document.getElementById('speech').onclick = function () {
  if (this.innerHTML != 'x0' && me.innerHTML == stage.innerHTML) {
    swal('選擇一個威脅', {
      buttons: {
        cancel: 'Cancel!',
        take0: { text: '雨天', value: 'Rain' },
        take1: { text: '雪天', value: 'Snow' },
        take2: { text: '夜晚', value: 'Night' },
        take3: { text: '子彈', value: 'Bullet' },
        take4: { text: '哨子', value: 'Whistle' },
        take5: { text: '面具', value: 'Mask' },
      },
    }).then((value) => {
      if (value != null) {
        let speech = '{"order":"speech","choose":"' + value + '"}';
        w.send(speech);
      }
    });
  }
};

players = getElementsByIdStartsWith('players', 'span', 'player');
threats = getElementsByIdStartsWith('players', 'div', 'threat');

///////////////////////////////////////////////////////////////////

w = new WebSocket('ws://' + HOST + '/my_endpoint');
/*
w.onopen = function () {
	console.log("Websocket connection enstablished");
};

w.onclose = function () {
	appendMessage("<div><center><h3>Disconnected</h3></center></div>");
};
*/
w.onmessage = function (message) {
  var jsonArray = JSON.parse(message.data);
  console.log(jsonArray);

  if (jsonArray['Process'] == 'game') {
    for (let i = 0; i < 7; i++) {
      if (i < jsonArray['NoMansLand'].length) {
        lands[
          i
        ].style = `background-image:url(./card/${jsonArray['NoMansLand'][i]['Name']}.png);`;
      } else {
        lands[i].style.backgroundImage = '';
      }
    }
    // 清空threat
    for (let i = 0; i < threats.length; i++) {
      var divs = threats[i].getElementsByTagName('div');
      while (divs.length > 0) {
        threats[i].removeChild(divs[0]);
      }
    }
    //放入threat
    for (let i = 0; i < jsonArray['Threats'].length; i++) {
      for (let j = 0; j < jsonArray['Threats'][i].length; j++) {
        var odiv = document.createElement('div');
        odiv.style = `background-image:url(./icon/${jsonArray['Threats'][i][j]}.png);`;
        threats[i].appendChild(odiv);
      }
    }

    for (let i = 0; i < jsonArray['Players'].length; i++) {
      players[i].innerHTML = jsonArray['Players'][i];
    }

    trialnum = jsonArray['TM'][0];
    moralenum = jsonArray['TM'][1];

    stage.innerHTML = jsonArray['Stage'];

    switch (stage.innerHTML) {
      case 'DrawCard': {
        swal({
          title: '遊戲就緒',
          icon: 'success',
          button: '抽牌',
        }).then(() => {
          w.send('{"order":"draw"}');
        });
        break;
      }
      case 'SupportSkip': {
        if (jsonArray['RoundWin']) {
          mission = '任務成功';
          missionicon = 'success';
        } else {
          mission = '任務失敗';
          missionicon = 'error';
        }

        swal({
          title: mission,
          icon: missionicon,
        }).then(() => {
          swal({
            title: '沒有人得到支援',
            icon: 'error',
            button: '確認',
          }).then((value) => {
            w.send('{"order":"supportEnd","choose":""}');
          });
        });
        break;
      }
      case 'Winner,Winner,Chicken Dinner': {
        swal({
          title: '遊戲勝利！',
          text: '和平再現!',
          icon: './card/Peace.png',
        });
        break;
      }
      case 'ID重複': {
        me.innerHTML = '';
        swal('ID重複!', '請重新登入', 'error');
        break;
      }
      case me.innerHTML: {
        swal('你的回合!', '請出牌,使用英雄能力,演說或撤退');
        break;
      }
    }
    switch (stage.innerHTML.split(':')[0]) {
      case 'Support': {
        //支援
        if (jsonArray['RoundWin']) {
          mission = '任務成功';
          missionicon = 'success';
        } else {
          mission = '任務失敗';
          missionicon = 'error';
        }

        if (jsonArray['Stage'].split(':')[1] == me.innerHTML) {
          swal({
            title: mission,
            icon: missionicon,
          }).then(() => {
            swal({
              title: '得到支援',
              icon: 'success',
              buttons: {
                take0: { text: '消除身上威脅', value: 'Threat' },
                take1: { text: '回覆英雄能力', value: 'Lucky' },
              },
            }).then((value) => {
              w.send('{"order":"supportEnd","choose":"' + value + '"}');
            });
          });
        } else {
          swal({
            title: mission,
            icon: missionicon,
          }).then(() => {
            swal({
              title: jsonArray['Stage'].split(':')[1] + '得到支援',
              icon: 'info',
              button: '確認',
            }).then(() => {
              w.send('{"order":"supportEnd","choose":""}');
            });
          });
        }
        break;
      }
      case 'Leader': {
        if (jsonArray['Stage'].split(':')[1] == me.innerHTML) {
          swal('選擇抽牌張數', {
            buttons: {
              take0: { text: '抽一張', value: '1' },
              take1: { text: '抽兩張', value: '2' },
              take2: { text: '抽三張', value: '3' },
              take3: { text: '抽四張', value: '4' },
            },
            closeOnConfirm: false,
          }).then((value) => {
            w.send('{"order":"newRound","choose":"' + value + '"}');
          });
        }
        break;
      }
      case '演說': {
        swal('演說階段', '目標:' + jsonArray['Stage'].split(':')[1]);
        break;
      }
      case 'Loser,Loser,now who’s dinner?': {
        swal({
          title: '遊戲失敗！',
          text: jsonArray['Stage'].split(':')[1],
          icon: './card/Monument.png',
        });
        break;
      }
    }
  }
  if (jsonArray['Process'] == 'hand') {
    for (let i = 0; i < 8; i++) {
      if (i < jsonArray['Handcard'].length) {
        hands[
          i
        ].style = `background-image:url(./card/${jsonArray['Handcard'][i]['Name']}.png);`;
      } else {
        hands[i].style = `background-image:url(./card/cardBack.png);`;
      }
    }
    hero.style = `background-image:url(./hero/${jsonArray['Hero']['Name']}.png)`;

    document.getElementById('speechNum').innerHTML =
      'x' + jsonArray['SpeechTime'];
    support0.innerHTML = 'x' + jsonArray['Support'][0];
    support1.innerHTML = 'x' + jsonArray['Support'][1];
    support2.innerHTML = 'x' + jsonArray['Support'][2];
    support3.innerHTML = 'x' + jsonArray['Support'][3];
  }
};

///////////////////////////////////////////////////
var trial = document.getElementById('trial');
var morale = document.getElementById('morale');
var restart = document.getElementById('restart');

restart.onclick = function () {
  w.send('{"order":"restart"}');
};

/* 牌庫檢查 */
var trialnum = 0;
trial.onclick = function () {
  swal({
    title: '考驗:' + trialnum,
  });
};

var moralenum = 0;
morale.onclick = function () {
  swal({
    title: '士氣:' + moralenum,
  });
};
/* 牌庫檢查結束 */
