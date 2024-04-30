<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function () {
    return {
      error: null,
      myID: localStorage.getItem("ID"),
      myUsername: localStorage.getItem("username"),
      userID: localStorage.getItem("userID"),
      isMyProfile: JSON.parse(localStorage.getItem("isMyProfile")),
      userProfile: {
        profile: {
          ID: 0,
          username: "",
          following_count: 0,
          follower_count: 0,
          post_count: 0
        },
        posts: [
          {
            ID: 0,
            like_count: 0,
            comment_count: 0,
            showCommentInput: false
          }
        ],
        followings: [],
        followers: [],
        banned_from: [],
      },
      profileOwner: "",
      newUsername: "",
      newPhoto: "",
      newDescription: "",
      newSearch: "",
      newComments: ["",],
      fullPost: {
        post: {
          ID: 0,
          like_count: 0,
          comment_count: 0,
        },
        like_owners: [
          {
            username: "",
          }
        ],
        full_comments: [
          {
            username: "",
            comment: {
              ID: 0,
              post_ID: 0,
              owner_ID: 0,
              owner_username: "",
              text: "",
              date_time: "",
            }
          }
        ],
      },
      showUsernameInput: false,
      showSearchInput: false,
      showUploadInput: false,
      showCommentInput: false,
      showLikeWindow: false,
      showCommentWindow: false,
      showLoading: false
    }
  },

  computed: {
    sortedPosts() {
      if (!this.userProfile.posts) {
        this.userProfile.posts = []
      }
      return this.userProfile.posts.slice().sort((a, b) => new Date(b.date_time) - new Date(a.date_time));
    },

    isFollowing() {
      if (!this.userProfile.followers) {
        this.userProfile.followers = []
      }
      return this.userProfile.followers.includes(Number(this.myID));
    },

    isBanned() {
      if (!this.userProfile.banned_from) {
        this.userProfile.banned_from = []
      }
      return this.userProfile.banned_from.includes(Number(this.myID));
    },
  },

  mounted() {
    this.getProfile()
  },

  methods: {
    async doLogout() {
      await localStorage.clear()
      this.$router.push({path: '/'})
    },

    async getProfile() {
      this.error = null
      this.showLoading = true;
        if (this.isMyProfile) {
        this.userProfile.profile.ID = this.myID
      } else {
        this.userProfile.profile.ID = this.userID
      }
      try {
        let response = await this.$axios.get(`/users/${this.userProfile.profile.ID}/profile`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile = response.data
        if (this.userProfile.posts == null) {
          this.userProfile.posts = []
        }
        if (this.userProfile.followings == null) {
          this.userProfile.followings = []
        }
        if (this.userProfile.followers == null) {
          this.userProfile.followers = []
        }
        if (this.userProfile.banned_from == null) {
          this.userProfile.banned_from = []
        }
        this.showLoading = false;
        await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
      } catch (e) {
        this.showLoading = false;
        if (e.response && e.response.status === 400) {
          this.error = "Failed to get user's profile";
        } else if (e.response && e.response.status === 404) {
          this.error = "User not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
      this.isMyProfile = true;
    },

    async searchUser() {
      this.error = null;
      this.showLoading = true;
      try {
        let response = await this.$axios.get(`/users?username=${this.profileOwner}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile.profile.ID = response.data
        try {
          let res = await this.$axios.get(`/users/${this.userProfile.profile.ID}/profile`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          })
          this.userProfile = res.data
          if (this.userProfile.posts == null) {
            this.userProfile.posts = []
          }
          if (this.userProfile.followings == null) {
            this.userProfile.followings = []
          }
          if (this.userProfile.followers == null) {
            this.userProfile.followers = []
          }
          if (this.userProfile.banned_from == null) {
            this.userProfile.banned_from = []
          }
          this.showLoading = false;
          this.isMyProfile = false;
          localStorage.setItem("isMyProfile", this.isMyProfile)
          await localStorage.setItem("userID", this.userProfile.profile.ID)
          this.$router.push({path: `/users/${this.userProfile.profile.ID}/profile`})
        } catch (e) {
          this.showLoading = false;
          await this.getProfile()
          if (e.response && e.response.status === 400) {
            this.error = "Failed to request user's profile";
          } else if (e.response && e.response.status === 404) {
            this.error = "User not found";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later";
          } else {
            this.error = e.toString();
          }
        }
      } catch (e) {
        await this.getProfile()
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request user's profile.";
        } else if (e.response && e.response.status === 404) {
          this.error = "User not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
        } else {
          this.error = e.toString();
        }
      }
      await this.toggleSearchInput();
      this.isMyProfile = true;
    },

    handleFileChange(event) {
      this.newPhoto = event.target.files[0];
      const reader = new FileReader();

      reader.readAsArrayBuffer(this.newPhoto);
    },

    async uploadPhoto() {
      try {
        let formData = new FormData();
        formData.append('file', this.newPhoto);
        formData.append('description', this.newDescription)
        await this.$axios.post("/posts", formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
            "Content-Type": "multipart/form-data"
          }
        })
        await this.togglePhotoInput()
        this.newDescription = "";
        await this.getProfile()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Description too long (max 35 characters)";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "File not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async deletePhoto(postID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        if (this.myUsername) {
          this.userProfile.posts = this.userProfile.posts.filter(post => post.ID !== postID);
          await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
          await this.getProfile()
        }
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to delete photo";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else {
          this.error = e.toString();
        }
      }
    },

    async setMyUserName() {
      if (this.newUsername === "") {
        this.error = "Username cannot be empty.";
      } else {
        this.error = null
        try {
          await this.$axios.put("/profile/setUserName", {username: this.newUsername}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          })
          this.userProfile.profile.username = this.newUsername;
          await localStorage.setItem("username", this.newUsername);
          await this.getProfile();
          window.location.reload();
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "This username is not available";
          } else if (e.response && e.response.status === 401) {
            this.error = "Operation not authorized";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later";
          } else {
            this.error = e.toString();
          }
        }
      }
    },

    async toggleLike(postID) {
      this.error = null;
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        const listOfLikeOwners = this.fullPost.like_owners || [];
        let isLiked = listOfLikeOwners.includes(this.myUsername);


        if (isLiked) {
          try {
            await this.$axios.delete(`/posts/${postID}/likes`, {
              headers: {
                Authorization: localStorage.getItem("username")
              }
            });
            this.userProfile.posts[i].like_count--;

          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to unlike";
            } else if (e.response && e.response.status === 401) {
              this.error = "Operation not authorized";
            } else if (e.response && e.response.status === 404) {
              this.error = "Post not found";
            } else {
              this.error = e.toString();
            }
          }

        } else {
          try {
            await this.$axios.put(`/posts/${postID}/likes`, {}, {
              headers: {
                Authorization: localStorage.getItem("username")
              }
            });
            this.userProfile.posts[i].like_count++;

          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to like";
            } else if (e.response && e.response.status === 401) {
              this.error = "Operation not authorized";
            } else if (e.response && e.response.status === 404) {
              this.error = "Post not found";
            } else {
              this.error = e.toString();
            }
          }

        }
        await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async showLikes(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        await localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showLikeWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async commentPhoto(postID) {
      this.error = null;
      try {
        let i = this.userProfile.posts.findIndex(post => post.ID === postID);
        let formData = new FormData();
        let tmp = this.newComments.reverse();
        formData.append('text', tmp[i])
        await this.$axios.post(`/posts/${postID}/comments`, formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
          }
        })
        this.userProfile.posts[i].comment_count++;
        await localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        await this.toggleCommentInput(postID)
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to comment";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async deleteComment(postID, commentID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}/comments/${commentID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        if (this.myUsername) {
          let i = this.userProfile.posts.findIndex(post => post.ID === postID);
          this.fullPost.full_comments = this.fullPost.full_comments.filter(comment => comment.post_ID !== comment.post_ID);
          this.userProfile.posts[i].comment_count--;
        }
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to delete comment";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Comment not found";
        } else {
          this.error = e.toString();
        }
      }
    },

    async showComments(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        await localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showCommentWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request comments";
        } else if (e.response && e.response.status === 401) {
          this.error = "Operation not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later";
        } else {
          this.error = e.toString();
        }
      }
    },

    async toggleFollow(UID) {
      if (this.userProfile.followers === null) {
        this.userProfile.followers = []
      }
      const isFollowed = this.userProfile.followers.includes(Number(this.myID))
      
      if (isFollowed) {
        try {
          await this.$axios.delete(`/users/${UID}/follow`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.profile.follower_count--;
          this.userProfile.followers = this.userProfile.followers.filter(user => user !== Number(this.myID))

        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to unfollow";
          } else if (e.response && e.response.status === 401) {
            this.error = "Operation not authorized";
          } else {
            this.error = e.toString();
          }
        }

      } else {
        try {
          await this.$axios.put(`/users/${UID}/follow`, {}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.profile.follower_count++;
          this.userProfile.followers.push(Number(this.myID))

        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to follow";
          } else if (e.response && e.response.status === 401) {
            this.error = "Operation not authorized";
          } else if (e.response && e.response.status === 404) {
            this.error = "User not found";
          } else {
            this.error = e.toString();
          }
        }
      }
      await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))

    },

    async toggleBan(UID) {
      const isBanned = this.userProfile.banned_from.includes(Number(this.myID))

      if (isBanned) {
        try {
          await this.$axios.delete(`/users/${UID}/ban`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.banned_from = this.userProfile.banned_from.filter(user => user !== Number(this.myID))

        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to unban";
          } else if (e.response && e.response.status === 401) {
            this.error = "Operation not authorized";
          } else {
            this.error = e.toString();
          }
        }

      } else {
        try {
          await this.$axios.put(`/users/${UID}/ban`, {}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.banned_from.push(Number(this.myID))

        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to ban";
          } else if (e.response && e.response.status === 401) {
            this.error = "Operation not authorized";
          } else if (e.response && e.response.status === 404) {
            this.error = "User not found";
          } else {
            this.error = e.toString();
          }
        }
      }
      await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
    },


    async closeLikeWindow() {
      this.showLikeWindow = false
    },

    async closeCommentWindow() {
      this.showCommentWindow = false
    },

    async toggleUsernameInput() {
      this.showUsernameInput = !this.showUsernameInput;
      if (!this.showUsernameInput) {
        await localStorage.removeItem(this.newUsername)
        this.newUsername = ""
      }
    },

    async toggleSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        await localStorage.removeItem(this.profileOwner)
        this.profileOwner = ""
      }
    },

    async togglePhotoInput() {
      this.showUploadInput = !this.showUploadInput;
      if (!this.showUploadInput) {
        await localStorage.removeItem(this.newPhoto)
        this.newPhoto = ""
        this.newDescription = ""
      }
    },

    async toggleCommentInput(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showCommentInput = !this.userProfile.posts[i].showCommentInput;
      if (!this.userProfile.posts[i].showCommentInput) {
        this.newComments[i] = ""
      }
      await localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
    },

  }
}
</script>

<template>

  <div class="d-flex position-relative">
    <div class="d-flex position-absolute top-0 end-0 mt-3">
      <ErrorMsg v-if="error" :msg="error"></ErrorMsg>
    </div>
  </div>
  <div>
    <div style="display: flex; justify-content: center;" v-if="this.userProfile.profile.username">
      <h2 class="h2" v-if="this.myUsername === this.userProfile.profile.username">Profile</h2>
      <h2 class="h2" v-else>Search</h2>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>


  </div>

  <div class="col-lg-8">
    <div >
      <div class="d-flex align-items-center">
        <h1 style="white-space: nowrap; margin-bottom: 0;">{{ userProfile.profile.username }}</h1>
        <h6 v-if="isBanned" style="color: red; margin-top: 20px">(banned)</h6>
      </div>
      <div class="col-lg-6 mx-5">
        <div class="d-flex mt-3 justify-content-between align-items-center">
          <h6 style="margin-right: 10px;">{{ userProfile.profile.follower_count }} Follower</h6>
          <h6 style="margin-right: 10px;">{{ userProfile.profile.following_count }} Following</h6>
          <h6>{{ userProfile.profile.post_count }} Post</h6>
          <div class="ms-5">
            <button type="button" class="btn" @click="toggleFollow(userProfile.profile.ID)" v-if="Number(this.myID) !== userProfile.profile.ID">
              <div style="display: flex; align-items: center;">
                <svg :class="{ 'feather': true, 'mb-2': true, 'me-1':true, 'is-following': isFollowing }">
                  <use :href="isFollowing ? '/feather-sprite-v4.29.0.svg#user-minus' : '/feather-sprite-v4.29.0.svg#user-plus'"/>
                </svg>
                <h6 class="btn-follow-text" v-if="isFollowing"> Unfollow </h6>
                <h6 class="btn-follow-text" v-else> Follow </h6>
              </div>
            </button>
            <button type="button" class="btn" @click="toggleBan(userProfile.profile.ID)" v-if="Number(this.myID) !== userProfile.profile.ID">
              <div style="display: flex; align-items: center;">
                <svg :class="{ 'feather': true, 'mb-2': true, 'me-1':true, 'is-following': isBanned }">
                  <use :href="isBanned ? '/feather-sprite-v4.29.0.svg#user-check' : '/feather-sprite-v4.29.0.svg#user-x'"/>
                </svg>
                <h6 class="btn-ban-text" v-if="isBanned">Unban</h6>
                <h6 class="btn-ban-text" v-else>Ban</h6>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
  </div>


  <div class="post-grid">
    <div class="loading-container" v-if="showLoading">
      <div class="loading">
        <h1 class="h1">Loading...</h1>
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#loader"/></svg>
      </div>
    </div>
    <div v-for="(post, index) in sortedPosts" :key="post.ID" class="post-container" v-else>
      <img v-if="post.file" :src="'data:image/jpeg;base64,' + post.file" alt="Post Image" class="post-image img-fluid align-content-center">
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <div class="d-flex justify-content-between">
        <p><span style="font-weight: bold;">{{ userProfile.profile.username }}</span>: {{ post.description }}</p>
        <p>{{post.date_time}}</p>
      </div>
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <button type="button" class="btn" @click="showLikes(post.ID)">
        Likes: {{ post.like_count}}
      </button>
      <button type="button" class="btn mb-1" @click="toggleLike(post.ID)">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
      </button>
      <div>
        <button type="button" class="btn" @click="showComments(post.ID)">
          Comments: {{ post.comment_count }}
        </button>
        <button type="button" class="btn mb-1" @click="toggleCommentInput(post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
        </button>
        <div v-if="post.showCommentInput" style="margin-right: 10px;">
          <input type="text" id="newComment" v-model="newComments[index]" class="form-control form-control-sm"
                 placeholder="comment text" aria-label="Recipient's comment" aria-describedby="basic-addon2">
          <button v-if="post.showCommentInput" type="button" class="btn btn-sm btn-primary" @click="commentPhoto(post.ID)">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
          </button>
        </div>
      </div>
      <div class="delete-button-container" v-if="this.myUsername === this.userProfile.profile.username">
        <button type="button" class="btn delete-button" @click="deletePhoto(post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
        </button>
      </div>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showLikeWindow">
    <div class="liked-users-modal">
      <h2>Likes</h2>
      <ul>
        <li v-for="username in this.fullPost.like_owners" :key="username">{{ username }}</li>
      </ul>
      <button @click="this.closeLikeWindow()">Close</button>
    </div>
  </div>

  <div class="liked-users-overlay" v-if="this.showCommentWindow">
    <div class="liked-users-modal">
      <h2>Comments</h2>
      <ul>
        <li v-for="fullComment in this.fullPost.full_comments" :key="fullComment.username">
          {{ fullComment.username + ": " + fullComment.comment.text }}
          <button v-if="fullComment.username === this.myUsername" type="button" class="btn delete-button mb-2" @click="deleteComment(fullComment.comment.post_ID, fullComment.comment.ID)">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
          </button>
        </li>
      </ul>
      <button @click="this.closeCommentWindow()">Close</button>
    </div>
  </div>


  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky d-flex flex-column">
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/stream" class="nav-link" style="font-size: 20px;">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#home"/>
                </svg>
                Home
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to="'/users/' + myID + '/profile'" class="nav-link"
                          style="font-size: 20px;" @click="getProfile">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#user"/>
                </svg>
                Profile
              </RouterLink>
            </li>
            <li class="nav-item mx-3" v-if="this.myUsername === this.userProfile.profile.username">
              <button type="button" class="btn btn-sm" style="font-size: 15px;" @click="togglePhotoInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#log-out"/>
                </svg>
                Upload
              </button>
              <div v-if="showUploadInput" style="margin-right: 10px;">
                <input type="file" id="newPhoto" @change="handleFileChange" class="form-control form-control-sm">
                <input type="text" id="newDescription" v-model="newDescription" class="form-control form-control-sm"
                       placeholder="Photo description (max 35 characters)" aria-label="Recipient's description"
                       aria-describedby="basic-addon2">
              </div>
              <button v-if="showUploadInput" type="button" class="btn btn-sm btn-primary" @click="uploadPhoto">Upload
              </button>
            </li>
            <li class="nav-item mx-3" v-if="this.myUsername === this.userProfile.profile.username">
              <button type="button" class="btn btn-sm" style="font-size: 15px;" @click="toggleUsernameInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#edit-3"/>
                </svg>
                Change Username
              </button>
              <div class="d-flex ">
                <input v-if="showUsernameInput" type="text" id="newUsername" v-model="newUsername"
                       class="form-control form-control-sm"
                       placeholder="New Username" aria-label="Recipient's username" aria-describedby="basic-addon2">
                <button v-if="showUsernameInput" type="button" class="btn btn-sm btn-primary ml-2 me-2"
                        @click="setMyUserName">
                  Change
                </button>
              </div>
            </li>
            <li class="nav-item mx-1">
                <button type="button" class="btn btn-sm" style="font-size: 20px;" @click="toggleSearchInput">
                  <svg class="feather mx-1">
                    <use href="/feather-sprite-v4.29.0.svg#search"/>
                  </svg>
                  <span style="font-weight: 500;">
                    Search
                  </span>
                </button>
              <div class="d-flex ">
                <input v-if="showSearchInput" type="text" id="Searched Username" v-model="profileOwner"
                       class="form-control form-control-sm"
                       placeholder="Who are you searching for?" aria-label="Recipient's username" aria-describedby="basic-addon2">
                <button v-if="showSearchInput" type="button" class="btn btn-sm btn-primary ml-2 me-2" @click="searchUser">
                  Search
                </button>
              </div>
            </li>
          </ul>
          <div class="mt-auto mb-3 mx-1">
            <RouterLink to="/" class="nav-link" style="font-size: 20px;" @click="doLogout">
              <svg class="feather mx-1">
                <use href="/feather-sprite-v4.29.0.svg#log-out"/>
              </svg>
              Logout
            </RouterLink>
          </div>
        </div>
      </nav>
    </div>
  </div>

</template>

<style>
.post-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(calc(33.33% - 15px), 1fr));
  grid-gap: 15px;
}

.post-container {
  position: relative;
  border: 2px solid #ccc;
  margin-bottom: 20px;
  width: 100%;
  margin-top: 50px;
  justify-content: center;
}

.post-image {
  object-fit: cover;
  aspect-ratio: 1.91/1;
  display: block;
  margin: 0 auto;
  height: 275px;
}


.delete-button-container {
  position: absolute;
  bottom: 0;
  right: 0;
}

.delete-button {
  font-size: 20px;
  color: red;
}

.liked-users-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.liked-users-modal {
  background-color: white;
  padding: 20px;
  border-color: #000000;
  border-width: 2px;
  border-style: solid;
  max-height: 80%;
  overflow-y: auto;
}

.followed .feather-user-minus {
  color: red;
}

.loading-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.5);
}

.loading {
  text-align: center;
}

.btn-follow-text {
  width: 60px;
  text-align: left;
}
.btn-ban-text {
  width: 40px;
  text-align: left;
}

</style>
