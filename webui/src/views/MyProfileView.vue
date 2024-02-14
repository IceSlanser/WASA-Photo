<script>
import {RouterLink} from "vue-router";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  components: {RouterLink},
  data: function () {
    return {
      error: null,
      myID: localStorage.getItem("ID"),
      myUsername: localStorage.getItem("username"),
      userProfile: JSON.parse(localStorage.getItem("userProfile")),
      isThisMyProfile: localStorage.getItem("isThisMyProfile"),
      profileOwner: "",
      newUsername: "",
      newPhoto: "",
      newText: "",
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
            }
          }
        ],
      },
      showUsernameInput: false,
      showSearchInput: false,
      showUploadInput: false,
      showCommentInput: false,
      showLikeWindow: false,
      showCommentWindow: false
    }
  },
  methods: {
    async doLogout() {
      localStorage.clear()
      this.$router.push({path: '/'})
    },

    async getProfile() {
      this.error = null
      try {
        let response = await this.$axios.get(`/users/${this.myID}/profile`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile = response.data
        localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        this.$router.push({path: `/users/${this.myID}/profile`})
      } catch (e) {
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
        formData.append('description', this.newText)
        await this.$axios.post("/posts", formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
            "Content-Type": "multipart/form-data"
          }
        })
        await this.getProfile();
        window.location.reload()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new username.";
        } else if (e.response && e.response.status === 401) {
          this.error = "uploadPhoto not authorized.";
        } else if (e.response && e.response.status === 404) {
          this.error = "Data not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
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
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        }
        window.location.reload()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new username.";
        } else if (e.response && e.response.status === 401) {
          this.error = "setMyUserName not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found";
        } else {
          this.error = e.toString();
        }
      }
    },

    async setMyUserName() {
      if (this.myUsername === "") {
        this.error = "Username cannot be empty.";
      } else {
        this.error = null
        try {
          await this.$axios.put("/profile/setUserName", {username: this.newUsername}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          })
          localStorage.setItem("username", this.newUsername);
          this.userProfile.profile.username = this.myUsername;
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
          window.location.reload()
        } catch (e) {
          if (e.response && e.response.status === 400) {
            this.error = "Failed to request new username.";
          } else if (e.response && e.response.status === 401) {
            this.error = "setMyUserName not authorized";
          } else if (e.response && e.response.status === 500) {
            this.error = "An internal error occurred, please try again later.";
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
        const likeOwnersArray = this.fullPost.like_owners || [];
        let isLiked = likeOwnersArray.includes(this.myUsername);


        if (isLiked) {
          await this.$axios.delete(`/posts/${postID}/likes`, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.posts[i].like_count--;

        } else {
          await this.$axios.put(`/posts/${postID}/likes`, {}, {
            headers: {
              Authorization: localStorage.getItem("username")
            }
          });
          this.userProfile.posts[i].like_count++;
        }
        localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
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
        localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showLikeWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async commentPhoto(postID) {
      try {
        let formData = new FormData();
        formData.append('text', this.newText)
        await this.$axios.post(`/posts/${postID}/comments`, formData, {
          headers: {
            Authorization: localStorage.getItem("username"),
          }
        })
        await this.getProfile();
        window.location.reload()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request new comment.";
        } else if (e.response && e.response.status === 401) {
          this.error = "commentPhoto not authorized.";
        } else if (e.response && e.response.status === 404) {
          this.error = "Data not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "An internal error occurred, please try again later.";
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
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        }
        window.location.reload()
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to delete.";
        } else if (e.response && e.response.status === 401) {
          this.error = "deleteComment not authorized";
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
        localStorage.setItem("fullPost", JSON.stringify(this.fullPost))
        this.showCommentWindow = true
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.error = "Failed to request post.";
        } else if (e.response && e.response.status === 401) {
          this.error = "toggleLike not authorized";
        } else if (e.response && e.response.status === 404) {
          this.error = "Post not found.";
        } else if (e.response && e.response.status === 500) {
          this.error = "Internal Server Error.";
        } else {
          this.error = e.toString();
        }
      }
    },

    async searchUser() {
      this.error = null
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
          localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
          this.$router.push({path: `/users/${this.myID}/profile`})
        } catch (e) {
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
      } catch (e) {
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
        localStorage.removeItem(this.newUsername)
        this.newUsername = ""
      }
    },

    async toggleSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        localStorage.removeItem(this.newUsername)
        this.newUsername = ""
      }
    },

    async togglePhotoInput() {
      this.showUploadInput = !this.showUploadInput;
      if (!this.showUploadInput) {
        localStorage.removeItem(this.newPhoto)
        this.newPhoto = ""
        this.newText = ""
      }
    },

    async toggleCommentInput(postID) {
      let i = this.userProfile.posts.findIndex(post => post.ID === postID);
      this.userProfile.posts[i].showCommentInput = !this.userProfile.posts[i].showCommentInput;
      if (!this.userProfile.posts[i].showCommentInput) {
        localStorage.removeItem((this.newText))
        this.newText = ""
      }
      localStorage.setItem("userProfile", JSON.stringify(this.userProfile))
    }
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
    <div style="display: flex; justify-content: center;">
      <h2 class="h2" v-if="this.isThisMyProfile">Profile</h2>
      <h2 class="h2" v-else>Search</h2>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>


  </div>

  <div class="col-lg-7">
    <div class="d-flex align-items-center">
      <h1 style="white-space: nowrap;">{{ userProfile.profile.username }}</h1>
      <div class="col-lg-6 mx-5">
        <div class="d-flex mt-3 justify-content-between">
          <h6 style="margin-right: 10px;">{{ userProfile.profile.follower_count }} Follower</h6>
          <h6 style="margin-right: 10px;">{{ userProfile.profile.following_count }} Following</h6>
          <h6>{{ userProfile.profile.post_count }} Post</h6>
        </div>
      </div>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
  </div>

  <div class="post-grid">
    <div v-for="post in userProfile.posts" :key="post.ID" class="post-container">
      <img :src="'data:image/jpeg;base64,' + post.file" alt="Post Image" class="post-image img-fluid">
      <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center mb-3 border-bottom"></div>
      <p><span style="font-weight: bold;">{{ userProfile.profile.username }}</span>: {{ post.description }}</p>
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
          <input type="text" id="newDescription" v-model="newText" class="form-control form-control-sm"
                 placeholder="comment text" aria-label="Recipient's comment" aria-describedby="basic-addon2">
          <button v-if="post.showCommentInput" type="button" class="btn btn-sm btn-primary" @click="commentPhoto(post.ID)">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
          </button>
        </div>
      </div>
      <div class="delete-button-container" v-if="this.isThisMyProfile">
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
            <li class="nav-item mx-3" v-if="isThisMyProfile">
              <button type="button" class="btn btn-sm" style="font-size: 15px;" @click="togglePhotoInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#log-out"/>
                </svg>
                Upload
              </button>
              <div v-if="showUploadInput" style="margin-right: 10px;">
                <input type="file" id="newPhoto" @change="handleFileChange" class="form-control form-control-sm">
                <input type="text" id="newDescription" v-model="newText" class="form-control form-control-sm"
                       placeholder="Photo description" aria-label="Recipient's description"
                       aria-describedby="basic-addon2">
              </div>
              <button v-if="showUploadInput" type="button" class="btn btn-sm btn-primary" @click="uploadPhoto">Upload
              </button>
            </li>
            <li class="nav-item mx-3" v-if="isThisMyProfile">
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
  padding: 10px;
  width: 100%;
  margin-top: 50px;
}

.post-image {
  margin-bottom: 10px;
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

.heart-icon {
  width: 15px;
  height: 15px;
  fill: red;
}

.keyboard-icon {
  width: 20px;
  height: 20px;
  fill: black;
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
</style>
