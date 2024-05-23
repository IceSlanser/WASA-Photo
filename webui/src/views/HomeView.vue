<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function() {
    return {
      error: null,
      myID: sessionStorage.getItem("ID"),
      myUsername: sessionStorage.getItem("username"),
      newUsername: "",
      profileOwner: "",
      newText: "",
      userID: 0,

      stream: [
        {
          post: {
            ID: 0,
            username: "",
            profile_ID: 0,
            file: "",
            description: "",
            like_count: 0,
            comment_count: 0,
            date_time: "",
            showCommentInput: false,
            showLikeWindow: false,
            showCommentWindow: false,
          },
          like_owners: [
            {
              owner_ID: 0,
              owner_name: ""
            }
          ],
          full_comments: [
            {
              username: "",
              comment: {
                ID: 0,
                post_ID: 0,
                owner_ID: 0,
                text: "",
                date_time:""
              }
            }
          ]
        }
      ],

      newComments: [],
      fullPost: {
        post: {
          ID: 0,
          like_count: 0,
          comment_count: 0,
        },
        like_owners: [
          {
            owner_ID: 0,
            owner_name: ""
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

      showSearchInput: false,
      showLoading: false,
      isMyProfile: false,
    }
  },

  mounted() {
    this.getStream();
  },

	methods: {
    async doLogout() {
      await sessionStorage.clear()
      this.$router.push({path: '/'})
    },

    async getStream() {
      this.error = null
      this.showLoading = true;
      try {
        let response = await this.$axios.get(`/stream`, {
          headers: {
            Authorization: sessionStorage.getItem("username")
          }
        })
        this.stream = response.data
        this.showLoading = false;
      } catch (e) {
        this.showLoading = false;
        if (e.response && e.response.status === 400) {
          this.error = "Failed to get stream";
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
    },

    async getProfile() {
      this.isMyProfile = true;
      sessionStorage.setItem("isMyProfile", JSON.stringify(this.isMyProfile));
      await this.$router.push({path: `/users/${this.myID}/profile`})
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
        this.userID = response.data
        this.showLoading = false;
        this.isMyProfile = false;
        sessionStorage.setItem("isMyProfile", this.isMyProfile)
        await sessionStorage.setItem("userID", this.userID)
        this.$router.push({path: `/users/${this.userID}/profile`})
      } catch (e) {
        await this.getStream()
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
      await this.toggleSearchInput()
    },

    async getUser(UID) {
      this.isMyProfile = false;
      sessionStorage.setItem("isMyProfile", this.isMyProfile)
      await sessionStorage.setItem("userID", UID)
      this.$router.push({path: `/users/${UID}/profile`})
    },

    async toggleLike(postID) {
      this.error = null;
      let i = this.stream.findIndex(post => post.post.ID === postID);
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
            this.stream[i].post.like_count--;
          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to delete.";
            } else if (e.response && e.response.status === 401) {
              this.error = "toggleLike not authorized";
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
            this.stream[i].post.like_count++;

          } catch (e) {
            if (e.response && e.response.status === 400) {
              this.error = "Failed to put.";
            } else if (e.response && e.response.status === 401) {
              this.error = "toggleLike not authorized";
            } else if (e.response && e.response.status === 404) {
              this.error = "Post not found";
            } else {
              this.error = e.toString();
            }
            setTimeout(() => {
              this.error = null;
            }, 3000);
          }
        }

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
        this.stream.forEach(post => {
          post.post.showLikeWindow = post.post.ID === postID;
        });
        this.stream.forEach(post => {
          post.post.showCommentWindow = false;
        });
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
    },

    async commentPhoto(postID) {
      this.error = null;
      try {
        let i = this.stream.findIndex(post => post.post.ID === postID);
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
        let j = this.stream.findIndex(post => post.post.ID === postID);
        this.stream[j].post.comment_count++;
        await this.toggleCommentInput(postID)
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
          let i = this.stream.findIndex(full_post => full_post.post.ID === postID);
          this.stream[i].post.comment_count--;
        }
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
        this.stream.forEach(post => {
          post.post.showLikeWindow = false;
        });
        this.stream.forEach(post => {
          post.post.showCommentWindow = post.post.ID === postID;
        });
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
        setTimeout(() => {
          this.error = null;
        }, 3000);
      }
    },

    async toggleSearchInput() {
      this.showSearchInput = !this.showSearchInput;
      if (!this.showSearchInput) {
        this.profileOwner = ""
      }
    },

    async toggleCommentInput(postID) {
      let i = this.stream.findIndex(post => post.post.ID === postID);
      this.stream[i].post.showCommentInput = !this.stream[i].post.showCommentInput;
      if (!this.stream[i].post.showCommentInput) {
        this.newComments = []
      }
    },

    async closeLikeWindow(postID) {
      let i = this.stream.findIndex(post => post.post.ID === postID);
      this.stream[i].post.showLikeWindow = false
    },

    async closeCommentWindow(postID) {
      let i = this.stream.findIndex(post => post.post.ID === postID);
      this.stream[i].post.showCommentWindow = false
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
    <div>
      <h2 class="h2" style="display: flex; justify-content: center">Home</h2>
    </div>
    <div class="mb-3 border-bottom"></div>
    <div v-if="!this.stream">
      <div class="h3" style="display: flex; justify-content: center; padding-top: 15rem">Follow someone to get started</div>
      <div class="h3" style="display: flex; justify-content: center">or</div>
      <div class="h3" style="display: flex; justify-content: center">post something yourself</div>
    </div>
  </div>



  <div class="post-grid" v-if="!showLoading">
    <div v-for="(post, index) in stream" :key="post.ID" class="post-container" >
      <img v-if="post.post.file" :src="'data:image/jpeg;base64,' + post.post.file" alt="Post Image" class="post-image img-fluid align-content-center">

      <div class="position-relative">
        <div class="d-flex justify-content-between pt-3 ">
          <p>
            <button class="btn no-border-btn no-padding-btn no-vertical-align-btn">
              <span class="username" @click="getUser(post.post.profile_ID)">
                {{ post.post.username }}:
              </span>
            </button>
            <span class="text">{{ post.post.description }}</span>
          </p>
          <p style="margin-right: 0.5rem; font-size: 0.8rem; font-style: italic">{{ post.post.date_time }}</p>
        </div>
        <div class="border-bottom"></div>
        <button type="button" class="btn no-vertical-align-btn" @click="showLikes(post.post.ID)">
          Likes: {{ post.post.like_count}}
        </button>
        <button type="button" class="btn no-vertical-align-btn" @click="toggleLike(post.post.ID)">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
        </button>
        <div style=" display: flex; padding-bottom: 0.35rem; padding-right: 0.35rem">
          <div style="display: inline-block;">
            <button type="button" class="btn no-vertical-align-btn" @click="showComments(post.post.ID)">
              Comments: {{ post.post.comment_count }}
            </button>
            <button type="button" class="btn no-vertical-align-btn" @click="toggleCommentInput(post.post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
            </button>
          </div>
          <div v-if="post.post.showCommentInput" style="display: flex; flex-grow: 1">
            <input type="text" id="newComment" v-model="newComments[index]" class="form-control form-control-sm" style="width: 100%"
                   placeholder="What do you want to comment?" aria-label="Recipient's comment" aria-describedby="basic-addon2">
            <button v-if="post.post.showCommentInput" type="button" class="btn no-vertical-align-btn btn-sm btn-primary" @click="commentPhoto(post.post.ID)">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
            </button>
          </div>
        </div>

        <div class="user-like-overlay" v-if="this.stream[index].post.showLikeWindow">
          <div class="user-like-modal">
            <div class="vertical-line">
              <svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
            </div>

            <div style="margin-left: 50%; transform: translateX(-50%); margin-top: 0.1rem;">
              <button class="btn no-border-btn no-padding-btn no-vertical-align-btn" style="display: flex; flex-direction: column; align-items: center">
                <span v-for="owner in this.fullPost.like_owners" :key="owner.owner_ID" class="like username btn no-vertical-align-btn no-border-btn" @click="getUser(owner.owner_ID)">
                  {{ owner.owner_name }}
                </span>
              </button>
            </div>

            <button class="btn close-button no-border-btn no-padding-btn no-vertical-align-btn" @click="this.closeLikeWindow(post.post.ID)">
              <svg class="feather" style="width: 1.5rem; height: 1.5rem"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
            </button>
          </div>
        </div>

        <div class="user-comment-overlay" v-if="this.stream[index].post.showCommentWindow">
          <div class="user-comment-modal">
            <div class="vertical-line">
              <svg class="feather icon"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
            </div>

            <ul style="margin-left: 1.5rem; margin-right: 1.70rem; margin-top: 0.5rem">
              <span v-for="fullComment in this.fullPost.full_comments" :key="fullComment.username" class="comment">
                <div class="d-flex">
                  <button class="btn no-border-btn no-padding-btn no-vertical-align-btn">
                    <div class="username btn no-border-btn no-padding-btn no-vertical-align-btn" @click="getUser(fullComment.comment.owner_ID)">
                      {{fullComment.username + ":  "}}
                    </div>
                  </button>
                  <div class="text">{{fullComment.comment.text }}</div>
                  <button v-if="fullComment.username === this.myUsername" type="button" class="btn delete-comment no-border-btn px-0"
                        @click="deleteComment(fullComment.comment.post_ID, fullComment.comment.ID)">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                  </button>
                </div>
                <div class="datetime">{{ fullComment.comment.date_time }}</div>
              </span>
            </ul>
            <button class="btn close-button no-border-btn no-padding-btn no-vertical-align-btn" @click="this.closeCommentWindow(post.post.ID)">
              <svg class="feather" style="width: 1.5rem; height: 1.5rem"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
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
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Home
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink :to="'/users/' + myID + '/profile'" class="nav-link" style="font-size: 20px;" @click="getProfile">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                Profile
              </RouterLink>
            </li>
            <li class="nav-item">
              <button type="button" class="btn no-border-btn" style="font-size: 20px;" @click="toggleSearchInput">
                <svg class="feather mx-1">
                  <use href="/feather-sprite-v4.29.0.svg#search"/>
                </svg>
                <span style="font-weight: 500;">Search</span>
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
  grid-template-columns: repeat(auto-fill, minmax(calc(33.33% - 15px), 1fr));
  grid-gap: 15px;
}

.post-container {
  position: relative;
  border: 2px solid #000;
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

.username {
  font-style: italic;
  font-weight: bold;
  font-size: 1rem;
  margin-left: .5rem
}

.text {
  font-size: 1rem;
  margin-left: .3rem
}

.loading-container {
  position: relative;
  top: 100%;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.user-like-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.user-like-modal {
  display: flex;
  flex-grow: 1;
  background-color: white;
  border-color: #000000;
  min-height: 100%;
  max-height: 100%;
  overflow-y: scroll;
}

.like {
  padding-bottom: .25rem;
  padding-right: .25rem;
}

.user-comment-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.user-comment-modal {
  display: flex;
  flex-grow: 1;
  background-color: white;
  border-color: #000000;
  min-height: 100%;
  max-height: 100%;
  overflow-y: scroll;
}

.comment {
  align-items: center;
  padding-bottom: .5rem;
}

.delete-comment {
  display: flex;
  height:100%;
  margin-left: .5rem;
  font-size: 1.1rem;
  color: red;
}

.close-button {
  position: absolute;
  inset-inline-end: 1.5rem;
  margin-top: .2rem;
  color: red;
}

.no-border-btn {
  border: transparent;
}

.no-padding-btn {
  padding: 0;
}

.no-vertical-align-btn {
  vertical-align: 0;
}

.vertical-line {
  position: absolute;
  border-right: 0.15rem solid #D3D3D3;
  height: 100%;
  padding-left: 3rem;
}

.icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  height: 2.5rem !important;
  width: 2.5rem !important;
}

</style>
