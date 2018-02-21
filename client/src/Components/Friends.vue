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
                      <i class="fa fa-circle online"></i> online
                      <!-- <i class="fa fa-circle offline"></i> offline -->
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
     // const id = this.$route.params.id;
      return axios.get(`http://localhost:8080/api/profile/2/friends`)
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
      width: 70px;
    & > img {
      width: 50px;
      height: 50px;
      border-radius: 50%;
    }
    .friend-data {
      float: left;
    }
  }
}

</style>
