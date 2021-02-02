import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import Icon from '@material-ui/core/Icon';
import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {Drawer} from '@material-ui/core';
import PageMenu from './PageMenu';

export interface ITopAppBarProps {
    title?: string;
    currentPage?: string;
    children?: React.ReactNode;
}

const TopAppBar: React.FC<ITopAppBarProps> = ({title, currentPage, children}) => {
    const menuAnchor = React.useRef<HTMLDivElement>(null);
    const [showMenu, setShowMenu] = React.useState(false);
    const toggleMenu = React.useCallback(() => setShowMenu(!showMenu), [showMenu]);
    const hideMenu = React.useCallback(() => setShowMenu(false), []);
    return <>
        <AppBar position='fixed'>
            <Toolbar variant='dense' disableGutters>
                <IconButton onClick={toggleMenu} color='inherit'><Icon>menu</Icon></IconButton>
                {title ? <Typography variant='h6'>{title}</Typography> : null}
                {children}
            </Toolbar>
            <div ref={menuAnchor} className='mdc-menu-surface--anchor' />
        </AppBar>
        <Drawer open={showMenu} onClose={hideMenu}>
            <Toolbar variant='dense' />
            <PageMenu currentPage={currentPage} />
        </Drawer>
    </>;
};

export default TopAppBar;
