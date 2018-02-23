"use strict"

<template>
    <div class="container">
        
        <header>
            <figure class="profile-banner">
                <img src="https://unsplash.it/975/300" alt="Profile banner"/>
            </figure>
            <figure class="profile-picture"
                    v-bind:style="{ 'background-image': 'url(' + profile.avatar_ref + ')' }">
            </figure>
            <div class="profile-stats">
                <ul>
                    <li><span>{{ profile.first_name }}  {{ profile.last_name }}</span></li>
                    <li>{{ profile.friends_num || 0 }}  friends</li>
                </ul>
                <button class="follow" :class="{followed: isFollowed}" v-on:click="follow">
                    {{ friendText }}
                </button>
            </div>
        </header>
        <body>
        <b-card>
            <b-media class="parent-post" v-for="post in posts" :key="post.post_id">
                <!-- <b-img slot="aside" :src="post.avatar_ref" height="64px" width="64" alt="placeholder" /> -->
                <h4>{{ post.post_title }}
                    <small> {{ post.post_last_update }}</small>
                </h4>
                <p>
                    {{ post.post_text }}
                </p>
                <img class="post-picture"v-if="post.file_link" :src="post.file_link">
                <!-- <div>:id="post.post_id"> -->
                <div @click="setLike(my_id, post.post_id, post.post_like)">
                    <!-- <div> Lol {{ post.post_like}}</div> -->
                    <img class="heart" v-show="checkLikedPost(my_id, post.post_like) > 0" src="../assets/like.png"/>
                    <img  class="heart" v-show="checkLikedPost(my_id, post.post_like) == 0" src="../assets/unlike.png"/>
                    <span>{{ post.post_like.length - 1 || 0}}</span>
                </div>
            </b-media>
        </b-card>
       
 

        </body>
    </div>
</template>

<script>


    const axios = require('axios');
    axios.defaults.withCredentials = true;

    const my_id = localStorage.getItem("id");

    export default {
        name: "profile",
        data: function () {
            return {
                isFollowed: false,
                friendText: "Add Friend"
            }
        },
        methods: {
            follow: function (event) {
                this.isFollowed = !this.isFollowed;
                if (this.isFollowed) {
                    this.friendText = "Unfollow";
                    this.friends = this.friends + 1;
                } else {
                    this.friends = this.friends - 1;
                    this.friendText = "Add Friend";
                }
                // console.log(this.friends);
            },

            checkLikedPost(myId, post_like) {
                // console.log("Post likes: " + post_like);
                if (post_like.length == 1)
                    return 0
                else if  (post_like.indexOf(myId) >= 0)
                    return 1
                else
                    return 0
                // console.log(post_likes.includes(myId));
            },
            setLike(my_id,post_id, post_like) {
               var i = post_like.indexOf(my_id);
                if(i != -1) {
                    post_like.splice(i, 1);
                }
                else {                  
                    post_like.push(my_id);
                }
                let request = "http://localhost:8080/api/post/";
                let appendRequest = request + p_id + "/like";
             let p_id = post_id;
                axios.put(appendRequest, post_like)
                .then(function (response) {
                  console.log(response.data);
                })
                .catch(function (error) {
                  error = true,
                  console.log(error.message);
            });

            }
        },
        asyncComputed: {
            profile: {
                get() {
                    let id = this.$route.params.id;
                return axios.get(`http://localhost:8080/api/profile/${id}`)
                    .then(resp => {
                        console.log(resp.data)
                        return resp.data
                    });
                    
                },
                default: {avatar_ref: 'https://unsplash.it/150/150'}
            },
            posts: {
                get() {
                    let id = this.$route.params.id;
                    let postList;
                    return axios.get(`http://localhost:8080/api/post/user/${id}`)
                        .then(resp => {
                            postList = resp.data;
                            // Some magic костыль
                            return Promise.all(postList.map(p => {
                                const i = p.user_id;
                                axios.get(`http://localhost:8080/api/profile/${i}`)
                                    .then(resp => {
                                        p["user_name"] = resp.data.first_name;
                                        p["post_last_update"] = p.post_last_update.slice(0, 16);
                                        p["avatar_ref"] = resp.data.avatar_ref;
                                    });
                            }))
                        }).then(() => postList);
                },
                default: {post_title: "Loading...", post_text: "Loading...", user_name: "Loading...", post_last_update: "Loading..."}
            }
        }
    }
</script>

<style scoped lang="less">
.parent-post {
    border: 1px solid lightgray;
    padding: 10px;
    margin-bottom: 10px;
    .post-picture{
        width: 150px;
        height: 150px;
    }
    
    .heart{
        width: 30px;
        height: 30px;

    }
}


.container {
    margin-top: 100px;
    padding: 0px;
}
    h5 {
        font-size: 30px;
    }

    h4 {
        display: block;
        font-size: 20px;
    }

    h4 > small {
        color: #aaaaaa;
        font-size: .5em;
        margin-left: 10px
    }

    header {
        box-shadow: 1px 1px 4px rgba(0, 0, 0, 0.5);
        margin: 25px auto 50px;
        height: 300px;
        position: relative;
        width: 100%;
        max-width: 1000px;
    }

    figure.profile-banner {
        left: 0;
        overflow: hidden;
        position: absolute;
        top: 0;
        z-index: 1;
        width: 100%;
    }

    figure.profile-picture {
        background-position: center center;
        background-size: cover;
        border: 5px #efefef solid;
        border-radius: 50%;
        bottom: -50px;
        box-shadow: inset 1px 1px 3px rgba(0, 0, 0, 0.2), 1px 1px 4px rgba(0, 0, 0, 0.3);
        height: 148px;
        left: 10%;
        position: absolute;
        width: 148px;
        z-index: 3;
    }

    div.profile-stats {
        bottom: 0;
        border-top: 1px solid rgba(0, 0, 0, 0.5);
        left: 0;
        padding: 15px 15px 15px 250px;
        position: absolute;
        right: 0;
        z-index: 2;

        /* Generated Gradient */
        background: -moz-linear-gradient(top, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.51) 3%, rgba(0, 0, 0, 0.75) 61%, rgba(0, 0, 0, 0.5) 100%);
        background: -webkit-gradient(linear, left top, left bottom, color-stop(0%, rgba(255, 255, 255, 0.5)), color-stop(3%, rgba(0, 0, 0, 0.51)), color-stop(61%, rgba(0, 0, 0, 0.75)), color-stop(100%, rgba(0, 0, 0, 0.5)));
        background: -webkit-linear-gradient(top, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.51) 3%, rgba(0, 0, 0, 0.75) 61%, rgba(0, 0, 0, 0.5) 100%);
        background: -o-linear-gradient(top, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.51) 3%, rgba(0, 0, 0, 0.75) 61%, rgba(0, 0, 0, 0.5) 100%);
        background: -ms-linear-gradient(top, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.51) 3%, rgba(0, 0, 0, 0.75) 61%, rgba(0, 0, 0, 0.5) 100%);
        background: linear-gradient(to bottom, rgba(255, 255, 255, 0.5) 0%, rgba(0, 0, 0, 0.51) 3%, rgba(0, 0, 0, 0.75) 61%, rgba(0, 0, 0, 0.5) 100%);
        filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#80ffffff', endColorstr='#80000000', GradientType=0);

    }

    div.profile-stats ul {
        list-style: none;
        margin: 0;
        padding: 0;
    }

    div.profile-stats ul li {
        color: #efefef;
        display: block;
        float: left;
        font-size: 16px;
        font-weight: normal;
        margin-right: 50px;
        text-shadow: 1px 1px 2px rgba(0,0,0,0.7)
    }

    div.profile-stats li span {
        display: block;
        font-size: 24px;
    }

    div.profile-stats button.follow {
        display: block;
        float: right;
        color: #ffffff;
        margin-top: 5px;
        text-decoration: none;

        /* This is a copy and paste from Bootstrap */
        background-color: #49afcd;
        text-shadow: 0 -1px 0 rgba(0, 0, 0, 0.25);
        background-color: #49afcd;
        background-image: -moz-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: -webkit-gradient(linear, 0 0, 0 100%, from(#5bc0de), to(#2f96b4));
        background-image: -webkit-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: -o-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: linear-gradient(to bottom, #5bc0de, #2f96b4);
        background-repeat: repeat-x;
        border-color: #2f96b4 #2f96b4 #1f6377;
        border-color: rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.25);
        filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#ff5bc0de', endColorstr='#ff2f96b4', GradientType=0);
        filter: progid:DXImageTransform.Microsoft.gradient(enabled=false);
        display: inline-block;
        padding: 4px 12px;
        margin-bottom: 0;
        font-size: 14px;
        line-height: 20px;
        box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.05);
        border-radius: 4px;
        cursor: pointer;
    }

    div.profile-stats button.follow.followed {

        /* Once again copied from Boostrap */
        color: #ffffff;
        text-shadow: 0 -1px 0 rgba(0, 0, 0, 0.25);
        background-color: #5bb75b;
        background-image: -moz-linear-gradient(top, #62c462, #51a351);
        background-image: -webkit-gradient(linear, 0 0, 0 100%, from(#62c462), to(#51a351));
        background-image: -webkit-linear-gradient(top, #62c462, #51a351);
        background-image: -o-linear-gradient(top, #62c462, #51a351);
        background-image: linear-gradient(to bottom, #62c462, #51a351);
        background-repeat: repeat-x;
        border-color: #51a351 #51a351 #387038;
        border-color: rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.25);
        filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#ff62c462', endColorstr='#ff51a351', GradientType=0);
        filter: progid:DXImageTransform.Microsoft.gradient(enabled=false);
    }

    header > h1 {
        bottom: -50px;
        color: #354B63;
        font-size: 40px;
        left: 350px;
        position: absolute;
        z-index: 5;
    }
</style>
