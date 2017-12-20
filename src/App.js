import React, { Component } from 'react';
import './App.css';

class SubscribeButton extends Component {
  constructor(props) {
    super(props);
  } 
  render() {
    return (
      <div onClick={this.props.subscribe} className='subscribeButton'>
        {this.props.isSubscribed? 'unsubscribe' : 'subscribe'}
      </div>
    );
  }
}

class Page extends Component {
  constructor(props) {
    super(props);
  } 
  getPage = (pageID) => {
    this.props.socket.send(JSON.stringify({
        "id": pageID,
        "name": "page.query",
        "args": {
            "pageId": pageID,
        }
    }));
    this.props.setPageDataLoading();
  };
  clickPage = () => {
    this.getPage(this.props.pageID);
    this.props.clickPage();
  }
  render() {
    var pageData = JSON.stringify(this.props.currentPageData);
    return (
      <div className='page' onClick={this.clickPage}>
        <div>
          <div className='pageTitle' onClick={this.clickPage}>{this.props.pageTitle}</div>
          <SubscribeButton subscribe={this.props.subscribe} isSubscribed={this.props.isSubscribed} />
        </div>
        {(this.props.showPage ? 
          <div className='pageData'>
             {(this.props.pageDataLoading ? 'page data loading' 
              : <div>{pageData}</div>
            )}
          </div>
          : null
        )}
      </div>
    );
  }
}

class Project extends Component {
  constructor(props) {
    super(props);
    this.state = {currentPage: null, subscribedPage: null};
  }  
  getPageList = (projectName) => {
    this.props.socket.send(JSON.stringify({
      "id": 0,
      "name": "page.list",
      "args": {
        "project": this.props.projectName,
      }
    }));
    this.props.setPagesLoading();
  }
  clickProject = () => {
    this.getPageList(this.props.projectName);
    this.props.clickProject();
  }
  clickPage = (pageid) => {
    if (this.state.currentPage !== pageid) {
      this.setState({currentPage: pageid});
    } else {
      this.setState({currentPage: null});
    } 
  }
  togglePageSubscribe = (pageID) => {
    if (this.state.subscribedPage === pageID) {
      this.setState({subscribedPage: null});
    } else {
      this.setState({subscribedPage: pageID});
    }
  }
  render() {
    var pages = [];
    if (this.props.pages) {
      for (var i = 0; i < this.props.pages.length; i++) {
        pages.push(
          <Page 
            key={i} 
            socket={this.props.socket}
            showPage={this.state.currentPage === this.props.pages[i].pageid}
            clickPage={this.clickPage.bind(this, this.props.pages[i].pageid)}
            pageTitle={this.props.pages[i].title}
            currentPageData={this.props.currentPageData} 
            pageID={this.props.pages[i].pageid} 
            pageDataLoading={this.props.pageDataLoading}
            setPageDataLoading={this.props.setPageDataLoading}
            subscribe={this.togglePageSubscribe.bind(this, this.props.pages[i].pageid)}
            isSubscribed={this.state.subscribedPage === this.props.pages[i].pageid}/>
        );
      }
    }

    return (
      <div className='project'> 
        <div>
          <div className='projectTitle' onClick={this.clickProject}>{this.props.projectName}</div>
          <SubscribeButton subscribe={this.props.subscribe} isSubscribed={this.props.isSubscribed}/>
        </div>
        {(this.props.showProject ? 
          <div className='pagesContainer'>
            {(this.props.pagesLoading ? 'pages loading' 
              : <div>{pages}</div>
            )}
          </div>
          : 
          null
        )}
      </div>
    );
  }
}

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      projects: [], 
      currentProject: null, 
      currentPages: [], 
      currentPageData: null,
      subscribedProject: null,
      subscribedPage: null,
      error: null,
      projectsLoading: false,
      pagesLoading: false,
      pageDataLoading: false,
    };
  }
  componentWillMount() {
    this.connection = new WebSocket('wss://wiki-meta-explorer.herokuapp.com/');

    this.connection.onopen = () => {
      if (this.connection.readyState) {
        this.getProjects(0);
      }
    };

    this.connection.onClose = (evt) => {   
      this.ws = new WebSocket('wss://wiki-meta-explorer.herokuapp.com/');
      //this.initWs();
    };

    this.connection.onError = (evt) => {
        this.setState({error: 'WebSocket error'});
    };

    this.connection.onmessage = (evt) => { 
      let resp = JSON.parse(evt.data);
      if (resp.name === 'project.list') {
        this.setState({
          projects: resp.data,
          projectsLoading: false,
        });
      } else if (resp.name === 'page.list') {
        this.setState({
          currentPages: resp.data,
          pagesLoading: false,
        })
      } else if (resp.name === 'page.query') {
        this.setState({
          currentPageData: resp.data,
          pageDataLoading: false,
        });
      } else if (resp.name === 'project.update') {
        this.setState({
          projects: resp.data,
        });
      } else if (resp.name === 'page.update') {
        this.setState({
          currentPages: resp.data,
        });
      }
    }; 
  }
  componentWillUnmount() {
    this.connection.close();
  }
  getProjects = (id) => {
    this.connection.send(JSON.stringify({
      "id": 0,
      "name": "project.list",
      "args": {}
    }));
    this.setState({projectsLoading: true});
  }
  showProject = (projectName) => {
    if (this.state.currentProject !== projectName) {
      this.setState({currentProject: projectName});
    } else {
      this.setState({currentProject: null, currentPages: null});
    }    
  }
  setPagesLoading = () => {
    this.setState({pagesLoading: true});
  }
  setPageDataLoading = () => {
    this.setState({pageDataLoading: true});
  }
  toggleProjectSubscribe = (projectName) => {
    if (this.state.subscribedProject === projectName) {
      this.setState({subscribedProject: null});
    } else {
      this.setState({subscribedProject: projectName});
    }
  }
  render() {
    var projects = [];
    for (var i = 0; i < this.state.projects.length; i++) {
      projects.push(
        <Project 
          key={i} 
          projectName={this.state.projects[i]} 
          clickProject={this.showProject.bind(this, this.state.projects[i])}
          showProject={this.state.currentProject === this.state.projects[i]} 
          pages={this.state.currentProject === this.state.projects[i] ? this.state.currentPages : null}
          currentPageData={this.state.currentPageData}
          socket={this.connection}
          pagesLoading={this.state.pagesLoading}
          setPagesLoading={this.setPagesLoading}
          pageDataLoading={this.state.pageDataLoading}
          setPageDataLoading={this.setPageDataLoading}
          isSubscribed={this.state.subscribedProject === this.state.projects[i]}
          subscribe={this.toggleProjectSubscribe.bind(this, this.state.projects[i])}/>
      );
    }
    return (
      <div>
        {this.state.error}
        {this.state.projectsLoading ? 
          'projects loading' : projects
        }
      </div>
    );
  }
}

export default App;
