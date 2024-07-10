describe('template spec', () => {
  it('Error handling', () => {
    cy.visit('/misc')
    cy.contains('Show error').click()
    cy.contains('new error').should('exist')
  })

  it('Panic handling', () => {
    cy.visit('/misc')
    cy.contains('Show panic').click()
    cy.contains('show panic').should('exist')
  })
})