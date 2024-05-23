<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function () {
    return {
      error: null,
      myID: sessionStorage.getItem("ID"),
      myUsername: sessionStorage.getItem("username"),
      userID: sessionStorage.getItem("userID"),
      isMyProfile: JSON.parse(sessionStorage.getItem("isMyProfile")),
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
            showCommentInput: false,
            showLikeWindow: false,
            showCommentWindow: false,
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
      newComments: [],
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
      showLoading: false
    }
  },

  computed: {
    isFollowing() {
      if (this.userProfile.followers == null || this.userProfile.followers === []) {
        return false
      }
      return this.userProfile.followers.includes(Number(this.myID));
    },

    isBanned() {
      if (this.userProfile.banned_from == null || this.userProfile.banned_from === []) {
        return false
      }
      return this.userProfile.banned_from.includes(Number(this.myID));
    },
  },

  mounted() {
    this.getProfile()
  },

  methods: {
    async doLogout() {
      await sessionStorage.clear()
      this.$router.push({path: '/'})
    },

    async getMyProfile() {
      this.isMyProfile = true;
      await this.getProfile()
    },

    async getUser(UID) {
      this.isMyProfile = false;
      this.userProfile.profile.ID = UID
      this.userID = UID
      await this.getProfile()
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
            Authorization: sessionStorage.getItem("username")
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
      this.isMyProfile = true;
    },

    async searchUser() {
      this.error = null;
      this.showLoading = true;
      try {
        let response = await this.$axios.get(`/users?username=${this.profileOwner}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        })
        this.userProfile.profile.ID = response.data
        try {
          let res = await this.$axios.get(`/users/${this.userProfile.profile.ID}/profile`, {
            headers: {
              Authorization: sessionStorage.getItem("username")
            }
          })
          this.userProfile = res.data
          this.showLoading = false;
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
      this.isMyProfile = false;
      await this.toggleSearchInput();
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
            Authorization: sessionStorage.getItem("username"),
            "Content-Type": "multipart/form-data"
          }
        })
        await this.togglePhotoInput()
        if (this.showUsernameInput) {
          await this.togglePhotoInput()
        }
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
    },

    async deletePhoto(postID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        })
        if (this.myUsername) {
          this.userProfile.posts = this.userProfile.posts.filter(post => post.ID !== postID);
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
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
              Authorization: sessionStorage.getItem("username")
            }
          })
          this.userProfile.profile.username = this.newUsername;
          this.myUsername = this.newUsername
          await sessionStorage.setItem("username", this.myUsername);
          if (this.showUploadInput) {
            await this.togglePhotoInput()
          }
          await this.toggleUsernameInput()
          await this.getProfile();
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
        }
      }
    },

    async toggleLike(postID) {
      this.error = null;
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        const listOfLikeOwners = this.fullPost.like_owners || [];
        let isLiked = listOfLikeOwners.some(owner => owner.owner_name === this.myUsername);


        if (isLiked) {
          try {
            await this.$axios.delete(`/posts/${postID}/likes`, {
              headers: {
                Authorization: sessionStorage.getItem("username")
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
            setTimeout(() => {
              this.error = null;
            }, 3000);
          }

        } else {
          try {
            await this.$axios.put(`/posts/${postID}/likes`, {}, {
              headers: {
                Authorization: sessionStorage.getItem("username")
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
    },

    async showLikes(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        this.userProfile.posts.forEach(post => {
          post.showLikeWindow = post.ID === postID;
        });
        this.userProfile.posts.forEach(post => {
          post.showCommentWindow = false;
        });
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
    },

    async commentPhoto(postID) {
      this.error = null;
      try {
        let i = this.userProfile.posts.findIndex(post => post.ID === postID);
        if (!this.newComments[i]) {
          this.newComments[i] = "";
        }
        let formData = new FormData();
        formData.append('text', this.newComments[i])

        await this.$axios.post(`/posts/${postID}/comments`, formData, {
          headers: {
            Authorization: sessionStorage.getItem("username"),
          }
        })
        let j = this.userProfile.posts.findIndex(post => post.ID === postID);
        this.userProfile.posts[j].comment_count++;
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }

    },

    async deleteComment(postID, commentID) {
      this.error = null
      try {
        await this.$axios.delete(`/posts/${postID}/comments/${commentID}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        })
        if (this.myUsername) {
          this.fullPost.full_comments = this.fullPost.full_comments.filter(full_comment => full_comment.comment.ID !== commentID);
          let i = this.userProfile.posts.findIndex(post => post.ID === postID);
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }

    },

    async showComments(postID) {
      this.error = null;
      try {
        let response = await this.$axios.get(`/posts/${postID}`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        });
        this.fullPost = response.data;
        let i = this.userProfile.posts.findIndex(post => post.ID === postID);
        this.userProfile.posts[i].showCommentWindow = true
        this.userProfile.posts.forEach(post => {
          post.showLikeWindow = false;
        });
        this.userProfile.posts.forEach(post => {
          post.showCommentWindow = post.ID === postID;
        });
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
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
              Authorization: sessionStorage.getItem("username")
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
        }

      } else {
        try {
          await this.$axios.put(`/users/${UID}/follow`, {}, {
            headers: {
              Authorization: sessionStorage.getItem("username")
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
        }
      }

    },

    async toggleBan(UID) {
      if (this.userProfile.banned_from == null) {
        this.userProfile.banned_from = []
      }
      const isBanned = this.userProfile.banned_from.includes(Number(this.myID))

      if (isBanned) {
        try {
          await this.$axios.delete(`/users/${UID}/ban`, {
            headers: {
              Authorization: sessionStorage.getItem("username")
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
        }

      } else {
        try {
          await this.$axios.put(`/users/${UID}/ban`, {}, {
            headers: {
              Authorization: sessionStorage.getItem("username")
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
          setTimeout(() => {
            this.error = null;
          }, 3000);
        }
      }
    },


    async closeLikeWindow(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showLikeWindow = false
    },

    async closeCommentWindow(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showCommentWindow = false
    },

    async togglePhotoInput() {
      this.showUploadInput = !this.showUploadInput;
      if (!this.showUploadInput) {
        this.newUsername = ""
        this.profileOwner = ""
        this.newPhoto = ""
        this.newDescription = ""
        this.newComments = []
        this.showSearchInput = false;
        this.showUsernameInput = false;
      }
    },

    async toggleUsernameInput() {
      this.showUsernameInput = !this.showUsernameInput;
      if (!this.showUsernameInput) {
        this.newUsername = ""
        this.profileOwner = ""
        this.newPhoto = ""
        this.newDescription = ""
        this.newComments = []
        this.showSearchInput = false;
        this.showUploadInput = false;
      }
    },

    async toggleSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        this.newUsername = ""
        this.profileOwner = ""
        this.newPhoto = ""
        this.newDescription = ""
        this.newComments = []
        this.showUploadInput = false;
        this.showUsernameInput = false;
      }
    },

    async toggleCommentInput(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showCommentInput = !this.userProfile.posts[i].showCommentInput;
      if (!this.userProfile.posts[i].showCommentInput) {
        this.newComments = []
      }
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



  <div class="loading-container" style="padding-top: 10rem" v-if="showLoading">
    <div style="text-align: center">
      <h1 class="h1">Loading...</h1>
      <div class="spinner-border"></div>
    </div >
  </div>

  <div v-else>
    <div style="display: flex; justify-content: center;">
      <h2 class="h2" v-if="this.myUsername === this.userProfile.profile.username">Profile</h2>
      <h2 class="h2" v-else>Search</h2>
    </div>
    <div class="mb-3 border-bottom"></div>
  </div>

  <div class="col-lg-10" v-if="!showLoading">
    <div >
      <div class="d-flex align-items-center">
        <h1 style="white-space: nowrap; font-size: 3rem">{{ userProfile.profile.username }}</h1>
        <h6 v-if="isBanned" style="color: red; margin-top: 2rem">(banned)</h6>
      </div>
      <div class="col-lg-5 mx-5">
        <div class="d-flex mt-2 justify-content-between align-items-center">
          <h6 style="margin-right: 10px;">{{ userProfile.profile.follower_count }} Follower</h6>
          <h6 style="margin-right: 10px;">{{ userProfile.profile.following_count }} Following</h6>
          <h6>{{ userProfile.profile.post_count }} Post</h6>
          <div class="ms-5">
            <button type="button" class="btn no-vertical-align-btn" @click="toggleFollow(userProfile.profile.ID)" v-if="Number(this.myID) !== userProfile.profile.ID">
              <div style="display: flex; align-items: center;">
                <svg :class="{ 'feather': true, 'mb-2': true, 'me-1':true, 'is-following': isFollowing }">
                  <use :href="isFollowing ? '/feather-sprite-v4.29.0.svg#user-minus' : '/feather-sprite-v4.29.0.svg#user-plus'"/>
                </svg>
                <h6 class="btn-follow-text no-vertical-align-btn" v-if="isFollowing"> Unfollow </h6>
                <h6 class="btn-follow-text no-vertical-align-btn" v-else> Follow </h6>
              </div>
            </button>
            <button type="button" class="btn no-vertical-align-btn" @click="toggleBan(userProfile.profile.ID)" v-if="Number(this.myID) !== userProfile.profile.ID">
              <div style="display: flex; align-items: center;">
                <svg :class="{ 'feather': true, 'mb-2': true, 'me-1':true, 'is-following': isBanned }">
                  <use :href="isBanned ? '/feather-sprite-v4.29.0.svg#user-check' : '/feather-sprite-v4.29.0.svg#user-x'"/>
                </svg>
                <h6 class="btn-ban-text no-vertical-align-btn" v-if="isBanned">Unban</h6>
                <h6 class="btn-ban-text no-vertical-align-btn" v-else>Ban</h6>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="this.isMyProfile" class=" col-lg-5 mb-3 border-bottom"></div>
    <div v-else class=" col-lg-6 mb-3 border-bottom"></div>
  </div>

  <div class="post-grid" v-if="!showLoading">
    <div v-for="(post, index) in userProfile.posts" :key="post.ID" class="post-container">
      <img v-if="post.file" :src="'data:image/jpeg;base64,' + post.file" alt="Post Image" class="post-image img-fluid align-content-center">
      <div class="position-relative">
        <div class="d-flex justify-content-between pt-3">
          <p>
            <button class="btn no-border-btn no-padding-btn no-vertical-align-btn">
              <span class="username btn no-border-btn no-padding-btn no-vertical-align-btn" @click="getUser(post.profile_ID)">
                {{ post.username }}:
              </span>
            </button>
            <span class="text">{{ post.description }}</span>
          </p>
          <p style="margin-right: 0.5rem; font-size: 0.8rem; font-style: italic">{{post.date_time}}</p>
        </div>
        <div class="border-bottom"></div>
        <button type="button" class="btn no-vertical-align-btn" @click="showLikes(post.ID)">
          Likes: {{ post.like_count}}
        </button>
        <button type="button" class="btn no-vertical-align-btn" @click="toggleLike(post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
        </button>
        <div style="display: flex; padding-bottom: 0.35rem; padding-right: 0.35rem">
          <div style="display: inline-block;">
            <button type="button" class="btn no-vertical-align-btn" @click="showComments(post.ID)">
              Comments: {{ post.comment_count }}
            </button>
            <button type="button" class="btn no-vertical-align-btn" @click="toggleCommentInput(post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
            </button>
          </div>
          <div v-if="post.showCommentInput" style=" display: flex; flex-grow: 1">
            <input type="text" id="newComment" v-model="newComments[index]" class="form-control form-control-sm" style="width: 100%"
                   placeholder="What do you want to comment?" aria-label="Recipient's comment" aria-describedby="basic-addon2">
            <button v-if="post.showCommentInput" type="button" class="btn btn-sm btn-primary no-vertical-align-btn" @click="commentPhoto(post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
            </button>
          </div>
          <div class="delete-button-container" v-if="this.myUsername === this.userProfile.profile.username">
            <button type="button" class="btn delete-photo no-border-btn no-vertical-align-btn" @click="deletePhoto(post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
            </button>
          </div>
        </div>

        <div class="user-like-overlay" v-if="this.userProfile.posts[index].showLikeWindow">
          <div class="user-like-modal">
            <div class="vertical-line">
              <svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
            </div>

            <div class="vertical-line"></div>

            <div style="margin-left: 4rem; margin-right: 1.70rem; margin-top: 0.1rem">
              <button class="btn no-border-btn no-padding-btn no-vertical-align-btn">
                <span v-for="owner in this.fullPost.like_owners" :key="owner.username" class="me-2 like username btn no-border-btn no-vertical-align-btn" @click="getUser(owner.owner_ID)">
                  {{ owner.owner_name }}
                </span>
              </button>
            </div>
            <button class="btn close-button no-border-btn no-padding-btn no-vertical-align-btn" @click="this.closeLikeWindow(post.ID)">
              <svg class="feather" style="width: 1.5rem; height: 1.5rem"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
            </button>
          </div>
        </div>

        <div class="user-comment-overlay" v-if="this.userProfile.posts[index].showCommentWindow">
          <div class="user-comment-modal">
            <div class="vertical-line">
              <svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
            </div>

            <ul style="margin-left: 1.5rem; margin-right: 1.70rem; margin-top: 0.5rem">
              <span v-for="fullComment in this.fullPost.full_comments" :key="fullComment.username" class="comment">
                <div class="d-flex">
                  <button class="btn no-border-btn no-padding-btn no-vertical-align-btn">
                    <div class="username btn no-padding-btn no-vertical-align-btn no-border-btn" @click="getUser(fullComment.comment.owner_ID)">
                     {{ fullComment.username + ":  " }}
                    </div>
                  </button>
                  <div class="text">{{ fullComment.comment.text }}</div>
                  <button v-if="fullComment.username === this.myUsername" type="button" class="btn delete-comment no-border-btn px-0" @click="deleteComment(fullComment.comment.post_ID, fullComment.comment.ID)">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                  </button>
                </div>
                <div class="datetime">{{ fullComment.comment.date_time }}</div>
              </span>
            </ul>
            <button class="btn close-button no-border-btn no-padding-btn no-vertical-align-btn" @click="this.closeCommentWindow(post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
            </button>
          </div>
        </div>

      </div>

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
              <RouterLink :to="'/users/' + this.myID + '/profile'" class="nav-link"
                          style="font-size: 20px;" @click="getMyProfile">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#user"/>
                </svg>
                Profile
              </RouterLink>
            </li>
            <li class="nav-item mx-3" v-if="this.myUsername === this.userProfile.profile.username">
              <button type="button" class="btn btn-sm no-vertical-align-btn" style="font-size: 15px;" @click="togglePhotoInput">
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
              <button v-if="showUploadInput" type="button" class="btn btn-sm btn-primary no-vertical-align-btn" @click="uploadPhoto">Upload
              </button>
            </li>
            <li class="nav-item mx-3" v-if="this.myUsername === this.userProfile.profile.username">
              <button type="button" class="btn btn-sm no-vertical-align-btn" style="font-size: 15px;" @click="toggleUsernameInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#edit-3"/>
                </svg>
                Change Username
              </button>
              <div>
                <input v-if="showUsernameInput" type="text" id="newUsername" v-model="newUsername"
                       class="form-control form-control-sm"
                       placeholder="New Username" aria-label="Recipient's username" aria-describedby="basic-addon2">
                <button v-if="showUsernameInput" type="button" class="btn btn-sm btn-primary ml-2 me-2 no-vertical-align-btn"
                        @click="setMyUserName">
                  Change
                </button>
              </div>
            </li>
            <li class="nav-item">
                <button type="button" class="btn no-border-btn" style="font-size: 20px;" @click="toggleSearchInput">
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
                <button v-if="showSearchInput" type="button" class="btn btn-sm btn-primary ml-2 me-2 no-vertical-align-btn" @click="searchUser">
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
  grid-template-columns: repeat(auto-fill, minmax(calc(33.33% - 1rem), 1fr));
  grid-gap: 0.5rem;
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
  top: 50%;
  transform: translateY(-50%);
  right: 0;
}

.delete-photo {
  font-size: 1.1rem;
  color: red;
}

.loading-container {
  position: relative;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.btn-follow-text {
  width: 60px;
  text-align: left;
}
.btn-ban-text {
  width: 40px;
  text-align: left;
}

.datetime {
  font-size: 0.75rem;
  padding-left: 1.25rem;
  padding-bottom: 0.75rem;
}

</style>
