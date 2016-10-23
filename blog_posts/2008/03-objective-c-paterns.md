Id: 1203
Title: Objective-C patterns
Tags: objective c,cocoa
Date: 2008-03-15T15:42:10-07:00
Format: Markdown
--------------
MVC:

-   model is the data. gives the data to view. view reads the data.
    controller changes data
-   view is the display. reads data from model. reports user interaction
    to controller
-   controller. changes the data based on user interaction in the view

Target/Action (command):

-   in cocoa, this is how responding to actions from controls is
    implemented
-   e.g. when button is pressed, framework sends a message to
    application code
-   general way in which object (e.g. button) can send messages to
    undetermined object (part of application code you write)
-   IB (Interface Builder) does that when you create new control, class
    and instance and ctrl-drag defines target action

Delegation:

-   messages send to Holder object are forwarded to Delegate
-   Holder needs an instance of Delegate and a way to set it

Chain of responsibility:

-   a message is handled by the most appropriate object out of many
-   used in cocoa for handling events. An event message is sent to a
    first responder object (e.g. if clicked on a button, button objects
    is the first responder). If the object doesn’t know how to handle
    the message, it sends it to the next object in chain etc.

