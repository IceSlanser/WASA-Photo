<script>
import {RouterLink} from "vue-router";

export default {
  components: {RouterLink},
  data: function() {
    return {
      error: null,
      userProfile: {
        Profile: {
          ID: localStorage.getItem("ID"),
          username: localStorage.getItem("username"),
          FollowingCount: 0,
          FollowerCount: 0,
          PostCount: 0
        },
        Posts: [
          {
            ID: 0,
            ProfileID: 0,
            File: [],
            Description: "",
            LikeCount: 0,
            CommentCount: 0,
            DateTime: Date
          }
        ],
        Followings: [
          {
            ID: 0
          }
        ],
        Followers: [
          {
            ID: 0
          }
        ]
      }
    }
  },
  methods: {
    async doLogout() {
      localStorage.clear()
      this.$router.push({ path : '/'})
    },
    async uploadPhoto() {

    },
    async getMyProfile() {
      try {
        let response = await this.$axios.get(`/users/${this.userProfile.Profile.ID}/profile`, {
          headers: {
            Authorization: localStorage.getItem("username")
          }
        })
        this.userProfile = response.data

        localStorage.setItem("userProfile", JSON.stringify(this.userProfile));
        this.$router.push({path: `/users/${this.userProfile.Profile.ID}/profile`})
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.err = "Failed to request user's profile.";
        } else if (e.response && e.response.status === 404) {
          this.err = "User not found.";
        } else if (e.response && e.response.status === 500) {
          this.err = "An internal error occurred, please try again later.";
        } else {
          this.err = e.toString();
        }

      }
    }
  },
}
</script>

<template>

  <div class="input-group-append">
    <button class="btn btn-outline-danger px-6 py-1" type="button" @click="doLogout">Logout</button>
    <button class="btn btn-outline-success px-6 py-1" @click="uploadPhoto">Upload</button>
    <button class="btn btn-outline-dark px-6 py-1" @click="getMyProfile">Profile</button>
  </div>



  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
        <div class="position-sticky pt-3 sidebar-sticky">
          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>General</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
                Home
              </RouterLink>
            </li>
          </ul>

          <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>Secondary menu</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
                Item 1
              </RouterLink>
            </li>
          </ul>
        </div>
      </nav>
    </div>
  </div>

  <div>
    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Profile page</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
            Refresh
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
            Export
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
            New
          </button>
        </div>
      </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>

</template>

<style>
</style>
