import "bulma/css/bulma.css";
import "mapbox-gl/dist/mapbox-gl.css";
import React from "react";
import { MdPower } from "react-icons/md";
import ReactMapGL, { Marker, NavigationControl } from "react-map-gl";
import "./App.css";
import BGImg from "./hero-bg.jpg";

export default class App extends React.Component {
    constructor(props) {
        super(props);
        this.emojiMap = {
            space: "🚀",
            sedan: "🚗",
            suv: "🚙",
            heli: "🚁",
            plane: "🛩",
            clown: "🤡",
            "": "🐶"
        };
        this.initialState = {
            setup: {
                name: "",
                email: "",
                driver: true,
                carType: ""
            },
            user: {
                name: null,
                email: null,
                lat: 0,
                lng: 0,
                carType: "sedan",
                driver: true
            },
            users: [],
            viewport: {
                width: "100vw",
                height: "100vh",
                latitude: 37.7577,
                longitude: -122.4376,
                zoom: 16
            },
            loaded: false,
            connected: false
        };
        this.state = Object.assign({}, this.initialState);
        this.mapRef = React.createRef();
    }
    _onViewportChange = viewport => {
        this.setState({ viewport });
    };
    handleCarTypSelect = event => {
        event.persist();
        this.setState(prevState => {
            const nextState = Object.assign({}, prevState.setup, { carType: event.target.value });
            return { setup: nextState };
        });
    };
    hadleSetupName = event => {
        event.persist();
        this.setState(prevState => {
            const nextState = Object.assign({}, prevState.setup, { name: event.target.value });
            return { setup: nextState };
        });
    };
    hadleSetupEmail = event => {
        event.persist();
        this.setState(prevState => {
            const nextState = Object.assign({}, prevState.setup, { email: event.target.value });
            return { setup: nextState };
        });
    };
    handleDriverStatus = event => {
        event.persist();
        this.setState(prevState => {
            const nextState = Object.assign({}, prevState.setup, { driver: !prevState.setup.driver });
            return { setup: nextState };
        });
    };

    createPerson = async () => {
        if ("geolocation" in navigator) {
            navigator.geolocation.getCurrentPosition(
                async ({ coords: { latitude, longitude } }) => {
                    const body = { ...this.state.setup, lat: latitude, lng: longitude };
                    const createPersonReq = await fetch("https://api.floqars.com/people", {
                        mode: "cors",
                        method: "POST",
                        body: JSON.stringify(body),
                        headers: {
                            "Content-Type": "application/json"
                        }
                    });
                    const res = await createPersonReq.json();
                    this.setState(
                        () => ({
                            user: res,
                            loaded: true
                        }),
                        () => window.localStorage.setItem("floqars-user", JSON.stringify(res))
                    );
                },
                null,
                { enableHighAccuracy: true }
            );
        } else {
            window.alert("geolocation must be enabled :/");
        }
    };

    logout = () => {
        window.localStorage.removeItem("floqars-user");
        this.setState(prevState => Object.assign({}, this.initialState, { connected: prevState.connected }));
    };

    sendLocation = () => {
        if (!this.state.connected) {
            return;
        }
        const u = Object.assign({}, this.state.user);
        this.socket.send(JSON.stringify({ ...u, action: "broadcast" }));
    };

    handleLocationUpdates(msgData) {
        if (!Array.isArray(msgData)) {
            return;
        }
        this.setState(({ users, loaded }) => {
            if (!loaded) {
                return;
            }
            const incomingUsers = JSON.parse(msgData);
            const nextUsers = Array.from(users);
            for (const u of incomingUsers) {
                if (!users.find(u => u.email === u.email)) {
                    u.icon = this.emojiMap[u.carType];
                    nextUsers.push(u);
                }
            }
            return { users: nextUsers };
        });
    }

    async componentDidMount() {
        try {
            const localUser = window.localStorage.getItem("floqars-user");
            if (localUser) {
                this.setState(() => ({
                    user: JSON.parse(localUser),
                    loaded: true
                }));
            }
            if (!this.socket) {
                this.socket = new WebSocket("wss://live.floqars.com");
            }
            this.socket.onopen = () => {
                this.setState(() => ({ connected: true }));
                this.sendLocation();
            };
            this.socket.onclose = () => {
                this.setState(() => ({ connected: false }));
            };
            this.socket.onmessage = msg => this.handleLocationUpdates(msg.data);

            if ("geolocation" in navigator) {
                this.geoWatcher = navigator.geolocation.watchPosition(
                    ({ coords: { latitude, longitude } }) => {
                        this.setState(prev => {
                            const nextUser = Object.assign({}, prev.user, {
                                lat: latitude,
                                lng: longitude
                            });
                            return {
                                viewport: Object.assign({}, prev.viewport, { latitude, longitude }),
                                user: nextUser
                            };
                        });
                        this.sendLocation();
                    },
                    null,
                    { enableHighAccuracy: true }
                );
            }
        } catch (err) {
            console.error(err);
        }
    }

    componentWillUnmount() {
        navigator.geolocation.clearWatch(this.geoWatcher);
    }

    renderNavBar() {
        return (
            <nav className="navbar is-fixed-top" role="navigation" aria-label="main navigation">
                <div className="navbar-brand">
                    <a className="navbar-item">
                        <h1 style={{ fontWeight: 800, fontSize: 26 }}>FLOQARS</h1>
                    </a>
                    <span className="navbar-item">
                        <small>we drive, you close</small>
                    </span>
                </div>

                <div className="navbar-end">
                    <span className="navbar-item">
                        <MdPower style={{ color: this.state.connected ? "green" : "red" }} />
                    </span>

                    {this.state.loaded &&
                        this.state.user && [
                            <div key="a" className="navbar-item">
                                <a className="navbar-item">
                                    <h2>{this.state.user.name}</h2>
                                </a>
                            </div>,
                            <div key="b" className="navbar-item">
                                <a onClick={this.logout} className="button is-danger is-inverted">
                                    logout
                                </a>
                            </div>
                        ]}
                </div>
            </nav>
        );
    }

    renderSetupForm() {
        return [
            <section key={1} className="hero is-transparent is-medium" style={{ minHeight: "90vh" }}>
                <div className="hero-body">
                    <div className="columns is-centered">
                        <div className="column is-4 has-background-white" style={{ borderRadius: 5 }}>
                            <div className="field">
                                <h1 className="title has-text-black">Get started today</h1>
                                <p>
                                    FloQars from FloQast. A revolutionary new ride-share platform for accountants, by
                                    accountants. Built on the Excel platform and ground-breaking AI, FloQars can help
                                    you decrease your close time by <em>at least</em> 99%. Close on the way to work, eat
                                    lunch, and go home. Get started today!
                                </p>
                            </div>
                            <div className="field">
                                <label className="label">Name</label>
                                <div className="control">
                                    <input
                                        onChange={this.hadleSetupName}
                                        value={this.state.setup.name}
                                        className="input"
                                        type="text"
                                        placeholder="e.g Alex Smith"
                                    />
                                </div>
                            </div>

                            <div className="field">
                                <label className="label">Email</label>
                                <div className="control">
                                    <input
                                        onChange={this.hadleSetupEmail}
                                        value={this.state.setup.email}
                                        className="input"
                                        type="email"
                                        placeholder="e.g. alexsmith@gmail.com"
                                    />
                                </div>
                            </div>
                            <div className="field">
                                <label className="checkbox">
                                    <input
                                        onChange={this.handleDriverStatus}
                                        checked={this.state.setup.driver}
                                        type="checkbox"
                                    />{" "}
                                    I'm interested in being a driver
                                </label>
                            </div>
                            {this.state.setup.driver && (
                                <div className="field">
                                    <label className="label">Vehicle type</label>
                                    <div className="select">
                                        <select value={this.state.setup.carType} onChange={this.handleCarTypSelect}>
                                            <option value="space">Space ship 🚀</option>
                                            <option value="sedan">Sedan 🚗</option>
                                            <option value="suv">SUV 🚙</option>
                                            <option value="heli">Helicopter 🚁</option>
                                            <option value="plane">Plane 🛩</option>
                                            <option value="clown">Clown Car 🤡</option>
                                        </select>
                                    </div>
                                </div>
                            )}
                            <div className="field">
                                <a onClick={this.createPerson} className="button is-primary is-fullwidth">
                                    Get Closing 🚘
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        ];
    }
    render() {
        return (
            <div className="App">
                <header className="App-header">{this.renderNavBar()}</header>
                {!this.state.loaded ? (
                    this.renderSetupForm()
                ) : (
                    <ReactMapGL
                        {...this.state.viewport}
                        mapboxApiAccessToken={process.env.REACT_APP_MAPTOKEN}
                        onViewportChange={this._onViewportChange}
                    >
                        {this.state.users.map(u => (
                            <Marker key={u.email} offsetLeft={-20} offsetTop={-10} latitude={u.lat} longitude={u.lng}>
                                <span>{u.icon}</span>
                            </Marker>
                        ))}
                        <Marker
                            key={"user"}
                            offsetLeft={-20}
                            offsetTop={-10}
                            latitude={this.state.user.lat}
                            longitude={this.state.user.lng}
                        >
                            <span>{this.emojiMap[this.state.user.carType]}</span>
                        </Marker>
                        <div style={{ position: "absolute", left: 0, bottom: 100 }}>
                            <NavigationControl />
                        </div>
                    </ReactMapGL>
                )}
                <img src={BGImg} className="bgImg" />
            </div>
        );
    }
}
