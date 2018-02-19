<template>
    <div class="container">
        <div class="people-list">
           <!-- <div class="search">
            <input type="text" placeholder="search" />
            <i class="fa fa-search"></i>
          </div> -->
          <ul class="list">
            <li class="chat-title">
              <img class="avatar" src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/195612/chat_avatar_01.jpg"/>                   
              <div class="about">
                <div class="name">
                  Tetyana Prysiazhniuk
                </div>
                <div class="status">
                  <i class="fa fa-circle online"></i> online
                </div>
              </div>
            </li>

             <li class="chat-title">
              <img class="avatar" src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/195612/chat_avatar_01.jpg"/>                   
              <div class="about">
                <div class="name">
                  Tetyana Prysiazhniuk
                </div>
                <div class="status">
                  <i class="fa fa-circle online"></i> online
                </div>
              </div>
            </li>
          </ul>
      </div>
      <div class="chat">
        <div class="chat-header">
              <img class="avatar" src="https://s3-us-west-2.amazonaws.com/s.cdpn.io/195612/chat_avatar_01.jpg"/>                   
              <div class="chat-about">
                <div class="chat-with">
                  Tetyana Prysiazhniuk
                </div>
                <div class="status">
                  <i class="fa fa-circle online"></i> online
                </div>
              </div>
        </div>
        <div class="chat-history">
          <ul>
            <li>
              <div class="message-data float-right">
                <span class="message-data-time">12.12.2017</span>
                <span class="message-data-name">Alina Biletska</span>
                <i class="fa fa-circle me"></i>
              </div>
              <div class="message other-message float-right">
                  Hi Vincent, how are you? How is the project coming along?
              </div>
            </li>
            <li>
              <div class="message-data float-left">
                <span class="message-data-name">Alina Biletska
                  <i class="fa fa-circle online"></i>
                </span>
                <span class="message-data-time">12.12.2017</span>
              </div>
              <div class="message my-message float-left">
                  Are we meeting today? Project has been already finished and I have results to show you.
              </div>

            </li>
          </ul>
        </div>

      </div>
     
  <link rel=“stylesheet” href=“https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css“/>
  
  </div>

</template>

<script>

export default {
  name: 'app',
  data () {
    return {
      
    }
  }
}
</script>

<style scoped>
.container {
    margin: 0 auto;
    margin-top: 110px;
    max-width: 1000px;
    width: 100%;
    background: #444753;
    border-radius: 5px;
    height: auto;
    min-height: 300px;
    padding: 0px;
    display: flex;

}
ul {
  list-style-type: none;
}
.people-list {
  /* float:left; */
  width: 400px;
}
/* .chat-list {
  margin: 10px 15px;
  padding: 5px;
  display: flex;
} */
.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}
.people-list ul {
  padding: 20px;
  height: 770px;
}
.people-list ul li {
  padding-bottom: 20px;
  height: 70px;
}
.people-list .about {
    padding-left: 8px;
    float: left;
    /* margin-top: 8px; */
}
.people-list .name{
    color: white;

}
.people-list .status {
    color: #92959E;
}
.people-list img {
    float: left;
}
.chat {
    width: 100%;
    background: #F2F5F8;
    border-top-right-radius: 5px;
    border-bottom-right-radius: 5px;
    color: rgb(67, 70, 81);
}
.chat-header {
    padding: 15px;
    border-bottom: 2px solid white;
    height: auto;
    width: 100%;
    overflow: hidden;
}
.chat .chat-header .chat-about {
    float: left;
    padding-left: 10px;
}
.chat .chat-header .chat-with {
    font-weight: bold;
    font-size: 16px;
}
.chat .chat-history {
    padding: 30px;
    overflow-y: scroll;
    height: auto;
}
.chat .chat-history .message-data {
    margin-bottom: 15px;
    float: right;
}
.chat .chat-history .message-data-time {
    color: #a8aab1;
    padding-left: 6px;
}
.chat .chat-history .other-message {
    background: #94C2ED;
}
.chat .chat-history .message {
    color: white;
    padding: 18px 20px;
    line-height: 26px;
    font-size: 16px;
    border-radius: 7px;
    margin-bottom: 30px;
    width: 90%;
    position: relative;
}
.float-right {
    float: right;
}
.float-left{
  float: left;
}
.chat .chat-history .my-message {
    background: #86BB71;
}




.chat .chat-header img {
    float: left;
}
.online, .offline, .me {
  margin-right: 3px;
  font-size: 10px;
}

.online {
  color: #86BB71;
}

.offline {
  color: #E38968;
}

.me {
  color: #94C2ED;
}

</style>
