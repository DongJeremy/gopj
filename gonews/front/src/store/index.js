import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import moment from 'moment'

Vue.use(Vuex)

Vue.prototype.$http = axios
//axios.defaults.baseURL = 'http://127.0.0.1:8081'
axios.defaults.xsrfHeaderName = 'X-CSRFToken'
axios.defaults.xsrfCookieName = 'csrftoken'
axios.defaults.withCredentials = false
axios.defaults.timeout = '5000'
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8'

export default new Vuex.Store({
  state: {
    page: 1,
    keyword: '',
    newsTotal: 0,
    newsPer: 10,
    newsItems: [],
    jobId: '',
    jobStatus: '',
    show: false,
  },
  mutations: {
    SET_PAGE(state, page) { // 设置页面
      state.page = page
    },
    SET_KEYWORD(state, keyword) { // 设置搜索关键词
      state.keyword = keyword
    },
    SET_NEWS(state) { // 设置新闻数据
      let page = state.page || 1
      let keyword = state.keyword
      var url = '/api/news?page=' + page + '&keyword=' + keyword
      axios.get(url).then(function(resp) {
        var items = resp.data.items
        var news = []
        for (let i = 0, len = items.length; i < len; i++) {
          var sNews = {
            id: items[i].id,
            title: items[i].title,
            url: items[i].link,
            ctime: moment(items[i].ctime, 'YYYYMMDD').endOf('day').fromNow(),
            weekday: moment(items[i].ctime).format('dddd').substr(0, 3),
            newDay: 1
          }
          if (i > 0 && items[i].ctime == items[i - 1].ctime) { // 合并同一天的新闻
            sNews.newDay = 0
          }
          news.push(sNews)
        }
        state.newsTotal = resp.data.total
        state.newsPer = resp.data.per
        state.newsItems = news
      }).catch(function() {
        console.log('error')
      })
    },
    PULL_REPO(state) {
      var url = '/api/pull'
      axios.get(url).then(function(resp) {
        state.jobId = resp.data.jobid
        state.jobStatus = resp.data.status
      }).catch(function() {
        console.log('error')
      })
    },
    SET_STATUS(state, status) { // 设置status
      state.jobStatus = status
    },
    SET_STATUS_QUERY(state) {
      var url = '/api/job/status?id=' + state.jobId
      console.log(url)
      axios.get(url).then(function(resp) {
        state.jobId = resp.data.jobid
        state.jobStatus = resp.data.status
      }).catch(function() {
        console.log('error')
      })
    },
    SET_SHOW(state, show) { // 设置loading
      state.show = show
    }
  }
})