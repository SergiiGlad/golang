<template>
    <div>
        <header>
            <div class="container">
            <img src="../assets/logo1.png">
            <button type="button">Logout</button>
            </div>
        </header>
        <div class="menu">
            <!-- <b-nav pills>
                <b-nav-item active><router-link :to="'/profile/:id'">My profile</router-link></b-nav-item>
                <b-nav-item><router-link :to="'/friends'">Friends</router-link></b-nav-item>
                <b-nav-item><router-link :to="'/messages'">Messages</router-link></b-nav-item>
                <b-nav-item><router-link :to="'/photos'">Photos</router-link></b-nav-item>
            </b-nav> -->
            <div class="container-2">
            <router-link :to="'/profile/:id'">My profile</router-link>
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
 
}
</script>

<style scoped>
.router-link{
    background-color: black;
}
header {
    background-color: rgb(245, 235, 235);
    box-shadow: 1px 1px 1px 1px gray ;
    height: 50px;
    position: fixed;
    top:0;
    width: 100%;
    z-index: 10;    
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
    border: 1px solid black;
    color: black;
    margin-top: 10px;
}
.menu {
    position: fixed;
    top: 50px;
    height: 40px;
    width: 100%;
    background-color: white;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
}
.container{
    max-width: 700px;
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
    .container-2 a {
        color: black;
    }   
</style>
