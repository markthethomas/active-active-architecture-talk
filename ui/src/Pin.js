import React from "react";
import { Popup } from "react-map-gl";

export default class Pin extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            show: false
        };
    }
    toggle = () => {
        console.log("toggle");
        this.setState(({ show }) => {
            console.log({ show: !show });
            return { show: !show };
        });
    };
    render() {
        return [
            <span key={"a"} onClick={this.toggle}>
                {this.props.icon}
            </span>,
            this.state.show && (
                <Popup
                    key={"b"}
                    dynamicPosition={true}
                    tipSize={5}
                    anchor="top"
                    longitude={this.props.lng}
                    latitude={this.props.lat}
                    closeOnClick={true}
                    onClose={() => this.toggle()}
                >
                    {this.props.icon}
                </Popup>
            )
        ];
    }
}
