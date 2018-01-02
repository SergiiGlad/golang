package main

import (
	"fmt"
	"os"
)

/* structure of a stack node */
type sNode struct {
	data rune
	next *sNode
}

/* Returns 1 if character1 and character2 are matching left
   and right Parenthesis */
func isMatchingPair(character1 rune, character2 rune) bool {

	if character1 == '(' && character2 == ')' {
		return true
	}

	if character1 == '{' && character2 == '}' {
		return true
	}

	if character1 == '[' && character2 == ']' {
		return true
	}

	return false

}

/*Return 1 if expression has balanced Parenthesis */
func areParenthesisBalanced(exp []rune) bool {

	i := 0

	/* Declare an empty character stack */
	var stack *sNode

	/* Traverse the given expression to check matching parenthesis */
	for i < len(exp) {
		/*If the exp[i] is a starting parenthesis then push it*/
		if exp[i] == '{' || exp[i] == '(' || exp[i] == '[' {

			push(&stack, exp[i])
		}

		/* If exp[i] is an ending parenthesis then pop from stack and
		   check if the popped parenthesis is a matching pair*/

		if exp[i] == '}' || exp[i] == ')' || exp[i] == ']' {

			/*If we see an ending parenthesis without a pair then return false*/
			if stack == nil {
				return false
			} else {

				if !isMatchingPair(pop(&stack), exp[i]) {
					return false
				}
			}

			/* Pop the top element from stack, if it is not a pair
			   parenthesis of character then there is a mismatch.
			   This happens for expressions like {(}) */

		}

		i++
	}

	/* If there is something left in expression then there is a starting
	   parenthesis without a closing parenthesis */
	if stack == nil {
		return true /*balanced*/
	} else {
		return false /*not balanced*/
	}
}

/* UTILITY FUNCTIONS */
/*driver program to test above functions*/
func main() {

	exp := []rune{'{', '(', ')', ')', '}', '[', ']'}

	fmt.Printf("%q\n", exp)

	if areParenthesisBalanced(exp) {
		fmt.Printf("Balanced\n ")
	} else {
		fmt.Printf("Not Balanced\n")
	}
	return
}

/* Function to push an item to stack*/
func push(top_ref **sNode, new_data rune) {
	/* allocate node */
	new_node := new(sNode)

	if new_node == nil {
		fmt.Printf("Stack overflow n")

		os.Exit(0)
	}

	/* put in the data  */
	new_node.data = new_data

	/* link the old list off the new node */
	new_node.next = (*top_ref)

	/* move the head to point to the new node */
	(*top_ref) = new_node
}

/* Function to pop an item from stack*/
func pop(top_ref **sNode) rune {
	var (
		res rune
		top *sNode
	)

	/*If stack is empty then error */
	if *top_ref == nil {
		fmt.Printf("Stack overflow n")

		os.Exit(0)
	} else {
		top = *top_ref
		res = top.data
		*top_ref = top.next

	}
	return res

}
