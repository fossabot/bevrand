import React, { Component } from 'react';
import Header from './Header';
import Randomizer from './Randomizer';
import Playlists from './PlaylistSelection';
import TopRolledBeverages from './TopRolledBeverages';
import AuthService from './AuthService';
import config from './ConfigService';

const Auth = new AuthService();

const getPlaylists = async (userName, isHomePage) => {
  // console.log('Received username: ', userName);

  let url = userName && !isHomePage ? `/api/playlists?username=${userName}` : '/api/v2/frontpage';
  let body;
  try {
    const response = await fetch(`${config.proxyHostname}${url}`);
    body = await response.json();
  } catch(err) {
    if(err.httpstatus === 404) {
      return null;
    }
    console.error('Error: ', err);
  }
  return body;
};

class RandomizeArea extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isLoading: true,
      isHomePage: this.props.history.location.pathname === '/'
    };

    this.changePlaylist = this.changePlaylist.bind(this);
  }

  componentDidMount() {
    let playlists;
    let userName;

    if(Auth.loggedIn()){
      // console.log('User is logged in')
      userName = Auth.getProfile().username;
    }

    getPlaylists(userName, this.state.isHomePage)
      .then(result => {
        playlists = result;
        if(!playlists) {
          return null;
        }
        // Previously this returned TGIF playlist, 
        // but due to the User page also using this,
        // Now the first playlist in the list is used as default
        return playlists[0];
      })
      .then(currentPlaylist => {
        this.setState({
          isLoading: false,
          playlists: playlists,
          currentPlaylist: currentPlaylist
        });
      })
      .catch(err => console.log(err.message));
  }

  changePlaylist(playlist) {
    this.setState({
      currentPlaylist: playlist
    });
  };

  render() {
    let userName = Auth.loggedIn() ? Auth.getProfile().username : null;
    let headerText;
    if (this.state.isLoading) {
      return (
        <div className="App">
          <p>App is still loading, please wait....</p>
        </div>
      )
    }

    if((!this.state.playlists || !this.state.currentPlaylist) && userName && !this.state.isHomePage) {
      headerText = `Welcome ${userName}, to your personal randomize area! `
        + 'Here you can create your own lists and randomize these lists. '
        + 'No personal playlists could be found, please add them first';
    } else if(userName && !this.state.isHomePage) {
      headerText = `Welcome ${userName}, to your personal randomize area! `
      + 'Here you can create your own lists and randomize these lists.';
    } else {
      headerText = 'All those Beverage choices bringing you down? Feeling lucky? Let fate quench your thirst with the Beverage Randomizer.'
        + 'Like you know, for randomizing your beverages.';
    }

    return (
      <div className="RandomizeArea">
      {/* Passing the props (of the main logged in user to the components like so:
        https://github.com/ReactTraining/react-router/issues/5521) */}
        <Header headerText={headerText}/>
        {this.state.currentPlaylist && <Randomizer userName={userName} playlist={this.state.currentPlaylist} />}
        {this.state.playlists && this.state.currentPlaylist && <Playlists playlists={this.state.playlists} onClick={this.changePlaylist} />}
        {this.state.currentPlaylist && <TopRolledBeverages playlist={this.state.currentPlaylist} />}
      </div>
    );
  }
}

export default RandomizeArea;
