import React from 'react';
import {shallow} from 'enzyme';
import TopAppBar from './TopAppBar';
import PageMenu from './PageMenu';
import {AppBar, IconButton, Menu, Toolbar, Typography} from '@material-ui/core';

describe('TopAppBar', () => {
    const title = 'Page Title';
    const menuItems = <PageMenu />;

    it('renders AppBar with title, menu and children', () => {
        const component = shallow(<TopAppBar title={title} menuItems={menuItems}><div id='children' /></TopAppBar>);

        expect(component.find(AppBar)).toHaveProp('position', 'fixed');
        expect(component.find(Toolbar).find(Typography)).toHaveText(title);
        expect(component.find(Toolbar).find('#children')).toExist();
        expect(component.find(Menu)).toHaveProps({open: false, children: menuItems});
    });
    it('renders AppBar without title', () => {
        const component = shallow(<TopAppBar menuItems={menuItems}><div id='children' /></TopAppBar>);

        expect(component.find(AppBar)).toHaveProp('position', 'fixed');
        expect(component.find(Toolbar).find(Typography)).not.toExist();
        expect(component.find(Toolbar).find('#children')).toExist();
        expect(component.find(Menu)).toHaveProps({open: false, children: menuItems});
    });
    it('displays menu when button is clicked', () => {
        const component = shallow(<TopAppBar title={title} menuItems={menuItems} />);

        component.find(Toolbar).find(IconButton).simulate('click');

        expect(component.find(Menu)).toHaveProp('open', true);
    });
});
