<template>
    <div>
        <header>
            <div class="container">
            <img src="../assets/logo1.png">
            <button type="button">Logout</button>
            </div>
        </header>
        <div class="menu">
            <div class="container-2">
                <router-link :to="{ name: 'Profile', params: {id: user_id}}">Profile</router-link>

            <!-- <router-link :to="'/profile/1'">My profile</router-link> -->
            <router-link :to="'/friends'">Friends</router-link>
            <router-link :to="'/messages'">Messages</router-link>
            <router-link :to="'/photos'">Photos</router-link>
            </div> 
        </div>
    </div>
</template>

<script>
export default {
  name: 'navigation',
  data () {
    return {

    }
  },
  computed:{
     user_id(){
         console.log(localStorage.getItem("id"));
          return localStorage.getItem("id");
      }
  }
 
}
</script>

<style scoped>
header {
    background-color: rgb(245, 235, 235);
    box-shadow: 0px 0px 15px -1px gray ;
    height: 50px;
    position: fixed;
    top:0;
    width: 100%;
    z-index: 21;    
}
img {
    height: 40px;
    width: auto;
    float: left;

}
button {
    float: right;
    border-radius: 4px;
    background-color: white;
    color: black;
    margin-top: 10px;
    outline: none;
    cursor: pointer;
}
.menu {
    position: fixed;
    top: 50px;
    height: 40px;
    width: 100%;
    background-color: #fafafa;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
    border-bottom: 1px gray;
    z-index: 20;
    }
.container{
    max-width: 1000px;
    width: auto;
}
.container-2 {
    margin-right: auto;
    margin-left: auto;
    padding-right: 15px;
    padding-left: 15px;
    width: 100%;
    max-width: 400px;
    display: flex;
    justify-content: space-around;
    margin: auto;
    }
a{
    text-decoration: none;
    color:black;
    cursor: pointer;
    text-shadow: black;
}

a.router-link-exact-active.router-link-active{
    font-weight: bold;
}
</style>
