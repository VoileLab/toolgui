describe('Input', () => {
  it('Textarea input', () => {
    cy.visit('http://localhost:3000/input')
    cy.get('textarea').type('testarea: 1')
    cy.get('textarea').blur()
    cy.get('textarea').type('2')
    cy.get('textarea').type('3')
    cy.get('textarea').blur()
    cy.contains('Value: testarea: 123')
  })

  it('Textbox input', () => {
    cy.visit('http://localhost:3000/input')
    cy.get('input[id=textbox_component_Textbox]').type('abc')
    cy.get('input[id=textbox_component_Textbox]').blur()
    cy.get('input[id=textbox_component_Textbox]').type('abc')
    cy.get('input[id=textbox_component_Textbox]').blur()
    cy.contains('Value: abcabc')
  })

  it('Checkbox', () => {
    cy.visit('http://localhost:3000/input')
    cy.get('input[type=checkbox]').click()
    cy.contains('Value: true').should('exist')
  })

  it('Button click', () => {
    cy.visit('http://localhost:3000/input')
    cy.contains('button').click()
    cy.contains('Value: true').should('exist')
  })

  it('Select', () => {
    cy.visit('http://localhost:3000/input')
    cy.contains('Value3').click()
    cy.contains('Value: Value3').should('exist')
  })
})