<template>
     <div class="container">
       <b-card no-body>
        <b-tabs card>
          <b-tab title="Friends" active>

             <div v-for="friend in friends" :key="friend.id">
              <b-row  class="friend-list">
                <div class="avatar" v-bind:style="{ 'background-image': 'url(' + friend.avatar_ref + ')' }"></div>
                <div class="friend-data">
                    <div class="name">
                      {{friend.first_name}} {{friend.last_name}}
                    </div>
                    <div class="status">
                      <!-- <i class="fa fa-circle online"></i> online -->
                       offline
                    </div>
                </div>
              </b-row>
            </div>

          </b-tab>
          <b-tab title="People" >
          </b-tab>
        </b-tabs>
      </b-card>

    </div>
</template>

<script>
const axios = require('axios');
export default {
  name: 'friends',
  data () {
    return {
    }
  },
  asyncComputed: {
    friends(){
      const id = localStorage.getItem("id");
      return axios.get(`http://localhost:8080/api/profile/${id}/friends`)
          .then(resp => {
            console.log(resp.data);
            return resp.data
          })
    }
  }
}
</script>

<style scoped lang="less">
.container {
  margin-top: 100px;
  padding: 0px;
}
.friend-list{
  padding: 15px 15px;
  border-bottom: 3px solid white;
  background: #F2F5F8;
    color: rgb(67, 70, 81);
  .avatar{
background-position: center center;
        background-size: cover;
        border: 3px #efefef solid;
        border-radius: 50%;
        bottom: -50px;
        box-shadow: inset 1px 1px 3px rgba(0, 0, 0, 0.2), 1px 1px 4px rgba(0, 0, 0, 0.3);
        height: 60px;
        left: 10%;
        // position: absolute;
        width: 60px;
        // z-index: 3;
        margin-right: 10px;
    }
    .friend-data {
      float: left;
    }
  }


</style>
