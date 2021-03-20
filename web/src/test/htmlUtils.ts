
export function createDiv() {
    const div = document.createElement('div');
    div.scrollTo = jest.fn();
    div.scrollIntoView = jest.fn();
    return div;
}
