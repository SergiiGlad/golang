<template>
     <div class="container">
      <div  v-for="friend in friends" :key="friend.id">
        <!-- {{friend}} -->
        <div class="col-3">
          <!-- <img @v-bind={{friend.avatar_ref}}> -->
        </div>
        <div class="col-9">
          <p>{{friend.first_name}} {{friend.last_name}}</p>
        </div>
      </div>
    </div> 
</template>

<script>
const axios = require('axios');
export default {
  name: 'app',
  data () {
    return {
      
    }
  },
  asyncComputed: {
    friends(){
     // const id = this.$route.params.id;
      return axios.get(`http://localhost:8080/profile/2/friends`)
          .then(resp => {
            console.log(resp.data);
            return resp.data
          })    
    }
  }
}

</script>

<style scoped>
.container {
  margin-top: 100px;
  min-height:200px ;
  background-color: white;
}


</style>
